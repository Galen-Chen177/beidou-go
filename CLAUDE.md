# CLAUDE.md

本文件为 Claude Code (claude.ai/code) 提供项目上下文，每次对话启动时自动加载。

## 项目概述

Beidou-Go 是 [BeiDouMS/BeiDou-Server](https://github.com/BeiDouMS/BeiDou-Server) 的 Go 语言迁移版，一个 **冒险岛 GMS v0.83** 服务端模拟器。原项目 Java 36% + JavaScript 61.5%，迁到 Go 后保留 JS 脚本（通过 goja 引擎执行）。

## 构建、测试、运行命令

```bash
# 构建
make build                          # go build -o bin/beidou-go ./cmd/server/

# 运行（先构建再启动）
make run                            # bin/beidou-go serve --config=config/config.yaml
make run CONFIG=config/config.yaml  # 指定配置文件

# 开发模式（文件监听自动重启，需要安装 air）
make dev

# 运行全部测试
make test                           # go test ./...

# 运行单个测试
go test ./internal/crypto/ -v -run TestDecryptCapturedPackets

# 仅执行数据库迁移
make build && ./bin/beidou-go migrate --config=config/config.yaml

# 下载依赖
make deps                           # go mod tidy && go mod download
```

## 架构

### 启动流程

```
cmd/server/main.go
  → cli.Run()                           # urfave/cli/v3 框架，serve / migrate 两个子命令
    → cli/serve.go:runServe()
      → config.Load()                   # yaml.v3 解析配置文件
      → store.InitDB() + AutoMigrate()  # gorm + MySQL，约 80 个 model
      → NewTCPServer(host)              # 单个 TCP 引擎，Login 和 Channel 共享
      → login.Server.Start()            # goroutine，默认端口 8484
      → channel.Server.Start()          # goroutine，默认端口 8485
      → 等待 SIGINT/SIGTERM → Shutdown()
```

Login 和 Channel 在同一个进程、同一个 `TCPServer` 实例上运行，通过不同端口区分。后续可以拆成独立进程。

### 网络协议

自定义 TCP 协议，**AES-256-OFB 流密码** + 外层自定义混淆：

1. **握手**：TCP 连接建立后，服务端立即明文发送 16 字节 `SERVER_HELLO`（包含 sendIV + recvIV）。客户端回复 `CLIENT_HELLO`（加密包，4B header + 2B body）。之后所有通信都是加密的。
2. **封包线上格式**：`[header: 4B][body: 加密后的 N 字节]`。4 字节 header 用 IV 编码（含版本号校验位和 XOR 混淆后的包体长度），不是简单的小端序长度。
3. **解密流程**（收包）：header 校验 → 解码包体长度 → AES-OFB 解密 → 自定义 `DecryptData` 反混淆 → 得到 `[opcode: 2B 小端序][data]`
4. **加密流程**（发包）：基于当前 sendIV 构造 header → 自定义 `EncryptData` 混淆 → AES-OFB 加密 → 写入 TCP

关键文件：`internal/crypto/crypto.go`（Crypto 接口）、`internal/crypto/maple_crypto.go`（AES-OFB + 混淆实现）、`internal/network/session.go`（ReadPacket/SendPacket）、`internal/network/codec/packet.go`（Packet 结构体）。

### 服务层 — 接口反置模式

为避免循环依赖，handler **接口**定义在 server 包中，**实现**放在 `handler/` 包中：

```
internal/server/login/login_handler.go    → LoginHandler 接口（7 个方法）
internal/server/channel/channel_handler.go → ChannelHandler 接口（1 个方法）
internal/server/handler/auth.go           → AuthHandler（实现 LoginHandler）
internal/server/handler/channel_handler.go → ChannelHandlerImpl（实现 ChannelHandler）
```

`cli/serve.go` 通过手动依赖注入把各组件组装起来（不用 DI 框架）。

### Login 服务器 (`internal/server/login/`)

Opcode 分发见 `login_server.go:dispatch()`：

| Opcode | 常量 | Handler 方法 | 说明 |
|--------|------|-------------|------|
| 0x01 | LoginCheckPassword | HandleCheckPassword | 账号密码验证 |
| 0x04/0x0B | LoginServerListRereq/Req | HandleServerList | 服务器列表 |
| 0x05 | LoginCharListReq | HandleCharList | 角色列表 |
| 0x06 | LoginServerStatusReq | HandleServerStatusRequest | 服务器状态 |
| 0x13 | LoginCharSelect | HandleCharSelect | 选择角色进入游戏 |
| 0x15 | LoginCheckCharName | HandleCheckCharName | 检查角色名 |
| 0x16 | LoginCharCreate | HandleCharCreate | 创建角色 |
| 0x23 | LoginClientHello | (无需处理) | 握手 IV 同步 |

### Channel 服务器 (`internal/server/channel/`)

| Opcode | 常量 | Handler 方法 | 说明 |
|--------|------|-------------|------|
| 0x14 | ChannelHello (PLAYER_LOGGEDIN) | HandlePlayerLoggedin | 角色进入频道 |
| 0x23 | LoginClientHello | (无需处理) | 握手 IV 同步 |

### 封包构造工具 (`internal/server/server_lib/`)

`packet.go` 包含所有服务端→客户端的封包构造函数，对齐 Java `PacketCreator.getCharInfo()`：
- `LoginStatusSuccess` / `LoginStatusFailed` — 登录成功/失败响应
- `ServerListEntry`、`EndOfServerList`、`LastConnectedWorld`、`RecommendedWorlds` — 服务器列表（包 6/7/8/9）
- `CharList`、`CharNameResponse`、`AddNewCharEntry` — 角色管理
- `ServerIP` — 选择角色后返回频道 IP 端口
- `AfterLoginError` — 登录错误响应
- `GetCharInfo` (SET_FIELD 0x7D) — 角色进入游戏的最关键封包，包含 10+ 个子封包

`password.go` — 密码处理：兼容 bcrypt / SHA-1 / SHA-512 三种 hash 格式，支持自动迁移到 bcrypt。

`world.go` — `WorldDataProvider`，从 `game_config` 表加载世界/频道配置，无配置时提供默认值。

### 数据层

- **Model** (`internal/model/`): 约 80 个 gorm 结构体，对应 MySQL 表（Account、Character、Inventoryitem、Inventoryequipment、Skill、Quest*、Guild 等）
- **Store** (`internal/store/`): `AccountStore`、`CharacterStore`、`GameConfigStore` — 对 gorm 查询的薄封装
- **AutoMigrate**: 启动时自动执行（见 `store/db.go`），会创建/更新所有 model 对应的表，不需要手动跑 SQL 迁移

### JS 脚本引擎 (`internal/script/`)

使用 **goja**（纯 Go 实现的 ES5.1，无 CGO 依赖）执行从原项目复用的 JS 脚本（NPC 对话、任务、事件、反应堆）。Runtime 通过 `sync.Pool` 复用。桥接对象（`cm`、`pi`、`ms`、`em`）通过 `vm.Set()` 注入 Go 侧实现。当前处于 **桩代码阶段** — `SetCM()`/`SetPI()` 还是 TODO。

### 尚未实现的部分

README 中的第三~七期系统基本未开始：
- `internal/service/` — 空目录（预留给业务逻辑层）
- `internal/api/` — gin 路由已定义，但 handler 全是桩代码
- `internal/script/` — 引擎骨架有了，桥接方法未实现
- Channel 端除了 0x14 和 0x23 之外的其他 opcode 均未处理
- 地图系统、聊天、战斗、NPC、背包、技能、任务等均未开始

### 测试

目前只有 `internal/crypto/crypto_test.go` — 用 **Wireshark 真实抓包数据**（SERVER_HELLO → CLIENT_HELLO → LOGIN 三包序列）验证加解密正确性。这是最关键的测试：加密层有 bug，上层协议全都会出问题。

## 配置文件

`config/config.yaml` — 通过 `yaml.v3` 加载，`config.Default()` 提供默认值。主要配置项：
- `login.port`（默认 8484）、`login.auto_register`（自动注册新账号）
- `channel.port`（默认 8485）、`channel.max_players`
- `database.*` — MySQL 连接参数
- `script.path`、`script.hot_reload`

## 项目约定

- **协议对齐**：从 Java 迁移某个系统时，在 Java 侧加日志打印 hex dump，Go 侧也打印，两相对照。封包格式是"真理"。
- **README 维护**：`internal/server/` 下每个子目录都有 README.md，用 ✅/⬜ 标记功能状态。每完成一个功能模块要同步更新对应的 README。
- **注释和日志**：可以用中文。README 是中文的。
- **设计文档**：`docs/` 目录下存放功能设计文档（如 `login-flow-analysis.md`、`0x14-channel-player-loggedin-plan.md` 等）。
- **参考 Java 项目**：本地路径 `D:\chj\project\golang\BeiDou-Server\`（Windows），`/home/chj/workspace/golang_galen/BeiDou-Server/`（WSL）。
- **每个系统迁移步骤**：读 Java 源码理解逻辑 → 定义 Go model → 实现 store/DAO → 实现业务逻辑 → 实现 handler 对接客户端 → 对接 JS 脚本（如果涉及）→ 客户端验证
