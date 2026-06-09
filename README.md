# Beidou-Go — 冒险岛 GMS v0.83 服务端 Go 迁移

> 从 [BeiDouMS/BeiDou-Server](https://github.com/BeiDouMS/BeiDou-Server) (Java) 迁移到 Go 语言。

---

## 目录

- [1. 项目背景与目标](#1-项目背景与目标)
- [2. 原始项目结构分析](#2-原始项目结构分析)
- [3. 技术栈选择](#3-技术栈选择)
- [4. 项目目录结构设计](#4-项目目录结构设计)
- [5. 五大核心问题及解决方案](#5-五大核心问题及解决方案)
  - [5.1 协议层 — TCP 通信与封包](#51-协议层--tcp-通信与封包)
  - [5.2 加解密层](#52-加解密层)
  - [5.3 游戏逻辑](#53-游戏逻辑)
  - [5.4 JS 脚本引擎](#54-js-脚本引擎)
  - [5.5 数据库设计](#55-数据库设计)
- [6. 分期实施计划](#6-分期实施计划)
- [7. 迁移策略](#7-迁移策略)
- [8. 开发环境与工具](#8-开发环境与工具)
- [9. 当前进度](#9-当前进度)
- [10. 待讨论问题](#10-待讨论问题)
- [11. 注意事项](#11-注意事项)

---

## 1. 项目背景与目标

### 源项目

| 项目 | BeiDou-Server |
|---|---|
| 仓库 | https://github.com/BeiDouMS/BeiDou-Server |
| 版本 | GMS v0.83 |
| 语言 | Java 36% + JavaScript 61.5% + Vue 1.7% + TypeScript 0.8% |
| 构建 | Maven (Java 21) |
| 数据库 | MySQL 8+ |
| 协议 | AGPL-3.0 |

### 迁移目标

- **语言**：Java → Go
- **JS 脚本**：保留 JS（Nashorn → goja）
- **管理后台**：Vue 前端暂时保留不变
- **数据库**：复用 MySQL 表结构
- **目标**：可编译、可运行、功能逐步对齐

---

## 2. 原始项目结构分析

```
BeiDou-Server/
├── gms-server/                  # Java 后端 (Maven)
│   ├── pom.xml
│   └── src/main/java/
│       ├── net/server/          # 网络层 (TCP Socket, Mina/Netty)
│       ├── client/              # 客户端会话、封包处理
│       ├── config/              # 配置 (YAML)
│       ├── constants/           # 常量、枚举、Opcode
│       ├── tools/
│       │   ├── data/            # 数据编解码
│       │   └── packets/         # 封包构造工具
│       ├── server/
│       │   ├── login/           # 登录服务器
│       │   ├── channel/         # 频道服务器
│       │   └── cashshop/        # 商城服务器
│       ├── scripting/           # Nashorn 脚本引擎封装
│       ├── provider/            # 数据提供者 (DAO 层)
│       └── service/             # 业务逻辑
├── gms-ui/                      # Vue.js 管理后台
└── scripts/                     # JS 脚本 (NPC/任务/事件)
```

### 服务端核心模块关系

```
Client (端口 8484) ──→ Login Server  ──→ 认证、角色列表
                              │
                              ▼
Client (端口 7575+) ──→ Channel Server ──→ 地图、战斗、社交、商城
                              │
                              ▼
                         MySQL 数据库
```

---

## 3. 技术栈选择

### 总览

| 层面 | Java (原) | Go (迁移后) | 状态 |
|---|---|---|---|
| 网络框架 | Netty / Mina | 标准库 `net` + `goroutine` | 无需第三方 |
| 序列化 | 手写 byte[] | 手写 `[]byte` + 工具函数 | 无需第三方 |
| 加密 | JCE (AES) | `crypto/aes` 标准库 | 无需第三方 |
| 数据库 ORM | JDBC / MyBatis | `gorm` | ✅ |
| JS 引擎 | Nashorn (JDK 内嵌) | `goja` | ✅ |
| CLI | Spring Boot CLI | `urfave/cli/v3` | ✅ |
| 配置文件 | YAML (Spring) | `gopkg.in/yaml.v3` 直接解析 | ✅ |
| 日志 | Log4j / Slf4j | `logrus` | ✅ |
| REST API | Spring Boot | `gin` | ✅ |
| 构建 | Maven | `go build` + `Makefile` | 无需第三方 |

---

### 3.1 CLI 框架 — `urfave/cli/v3`

| 方案 | 优点 | 缺点 | 选择？ |
|---|---|---|---|
| **urfave/cli/v3** | API 简洁，v3 已稳定（v3.8.0）；纯 flag 定义，无魔法；支持子命令、自动补全；MIT 协议 | 不支持配置文件解析（需配合 yaml 库） | ✅ |
| cobra | Kubernetes/Docker 等大项目在用，生态最大 | API 繁琐，需要分别定义 Command/Flag；和 viper 绑定深 | ❌ |
| stdlib `flag` | 零依赖 | 不支持子命令，不支持长短选项同时存在 | ❌ |

**选择理由**：`urfave/cli/v3` 够用且简洁。服务端只需要 `serve` 和 `migrate` 两个子命令，不需要 cobra 的重量。配置文件解析交给 `yaml.v3` 直接 Unmarshal，职责分离得更清楚。

---

### 3.2 配置文件 — `yaml.v3` 直接解析

| 方案 | 优点 | 缺点 | 选择？ |
|---|---|---|---|
| **yaml.v3 直接解析** | 极简：`yaml.Unmarshal(data, &cfg)` 一行搞定；无隐式行为，一切显式可控 | 不自动合并多源；无环境变量覆盖 | ✅ |
| viper | 多源合并（文件+环境变量+远程）；热加载 | 隐式行为太多，类型断言容易踩坑；配置优先级规则复杂，排查困难 | ❌ |
| envconfig | 从环境变量读取，适合 12-factor | 不适合结构化 YAML 配置（多层嵌套不好表达） | ❌ |

**选择理由**：游戏服务端配置通常是单个 YAML 文件（端口、数据库连接、游戏参数），不需要 viper 的多源合并和环境变量覆盖。`yaml.v3` 一行代码搞定，直观可控。

---

### 3.3 数据库 ORM — `gorm`

| 方案 | 优点 | 缺点 | 选择？ |
|---|---|---|---|
| **gorm** | Go 生态最流行（38k+ stars）；AutoMigrate、关联、事务、Hook 全包；文档和社区丰富 | 复杂查询时生成的 SQL 不够优；零值更新是常见坑；反射开销 | ✅ |
| sqlx | 接近原生 SQL，性能好；无反射开销 | 需要手写大量 SQL；无 AutoMigrate；关联查询要手动拼接 | ❌ |
| ent | Facebook 出品，类型安全；代码生成（schema→代码） | 学习曲线陡；生成代码量大；schema 改动需要重新生成 | ❌ |
| sqlc | 从 SQL 文件生成类型安全代码；接近原生性能 | 依赖手写 SQL 文件；动态查询支持弱 | ❌ |

**选择理由**：从 Java/MyBatis 迁移，gorm 是最接近传统 ORM 体验的选择。表结构复杂（角色、技能、背包、公会等几十张表），gorm 的 AutoMigrate 和关联加载能大幅减少 CRUD 代码。性能风险后续可以在热点路径用原生 SQL 兜底。

---

### 3.4 REST API 框架 — `gin`

| 方案 | 优点 | 缺点 | 选择？ |
|---|---|---|---|
| **gin** | Go 生态最流行（80k+ stars）；中间件生态丰富（cors/auth/logger）；渲染/绑定/校验开箱即用 | 路由不是 radix tree（性能略逊于 httprouter）；底层封装多 | ✅ |
| chi | 轻量，兼容 `net/http`；中间件链式调用优雅 | 生态比 gin 小；内置功能少（需要自己组装） | ❌ |
| echo | 类似 gin，性能稍好 | 社区比 gin 小；更新不如 gin 活跃 | ❌ |
| stdlib `net/http` | 零依赖 | 路由、参数绑定、校验全要手写 | ❌ |

**选择理由**：管理后台 API 是典型的 CRUD + Swagger 场景，gin 的 `ShouldBindJSON`、分组路由、中间件链最省事。Go 1.22 后 `net/http` 有改进，但 gin 的开箱体验仍然最好。

---

### 3.5 日志 — `logrus`

| 方案 | 优点 | 缺点 | 选择？ |
|---|---|---|---|
| **logrus** | 最老牌的 Go 日志库（25k+ stars）；结构化日志、Hook 机制、Formatter/Level 经典设计；社区资料最多 | 性能中等（有反射开销）；作者标记为"功能冻结"但仍在维护 | ✅ |
| slog | Go 1.21+ 标准库，零依赖；API 设计现代 | 生态还年轻，Hook 和 Formatter 不如 logrus 丰富；与 gin/gorm 集成不如 logrus 现成 | ❌ |
| zap | Uber 出品，性能最强（零反射）；结构化日志 | API 稍复杂（Sugar 包装层）；可读性不如 logrus | ❌ |
| zerolog | 零反射，性能接近 zap；API 更简洁 | 社区比 logrus 小很多 | ❌ |

**选择理由**：游戏服务端日志不是性能瓶颈。logrus 的 `WithFields` + `SetLevel` + `SetFormatter` 组合最直观，跟 gin、gorm 的集成方案也最成熟。slog 虽然零依赖，但当前生态还不如 logrus 丰富。

---

### 3.6 JS 脚本引擎 — `goja`

| 方案 | 优点 | 缺点 | 选择？ |
|---|---|---|---|
| **goja** | 纯 Go ES5.1，无 CGO；API 友好（`vm.Set()` 绑定 Go 函数）；成熟稳定；社区活跃 | 不支持 ES6+ 语法（let/const/箭头函数部分不支持）；性能不如 V8 | ✅ |
| v8go | 完整 V8 引擎（支持最新 ES）；性能最强 | 依赖 CGO + libv8（编译麻烦，二进制巨大）；API 不如 goja 友好 | ❌ |
| gopher-lua | Lua 5.1，纯 Go；游戏行业熟悉 | 需要重写全部 JS 脚本；社区不如 goja 活跃 | 以后候选 |

**选择理由**：原项目 61.5% 的代码是 JS 脚本，换语言成本太高。goja 纯 Go 编译无痛，桥接层实现简单。未来功能跑通后，可以评估迁移到 gopher-lua。

---

## 4. 项目目录结构设计

```
beidou_go/
├── cmd/
│   └── server/
│       └── main.go              # 入口，启动所有服务
├── internal/
│   ├── network/                 # 网络层
│   │   ├── tcp_server.go        # TCP 服务器封装
│   │   ├── session.go           # 客户端会话管理
│   │   └── codec/               # 封包编解码
│   │       ├── packet.go        # 封包结构体 (opcode + data)
│   │       ├── reader.go        # 封包读取器
│   │       └── writer.go        # 封包写入器
│   ├── crypto/                  # 加解密层
│   │   ├── aes.go               # AES/OFB 加密
│   │   └── shanda.go            # Shanda 自定义加密 (GMS 83)
│   ├── opcode/                  # 操作码定义
│   │   ├── login.go             # 登录相关 opcode
│   │   └── channel.go           # 频道相关 opcode
│   ├── model/                   # 数据模型 (对应数据库表)
│   │   ├── character.go
│   │   ├── account.go
│   │   ├── item.go
│   │   ├── skill.go
│   │   └── ...
│   ├── store/                   # 数据访问层 (DAO)
│   │   ├── db.go                # 数据库连接初始化
│   │   ├── account_store.go
│   │   ├── character_store.go
│   │   └── ...
│   ├── server/                  # 服务器实例
│   │   ├── login/               # 登录服务器
│   │   │   ├── login_server.go  # 启动、监听
│   │   │   └── handler/         # 登录封包处理器
│   │   │       ├── auth.go      # 认证处理
│   │   │       ├── world.go     # 世界/频道列表
│   │   │       └── character.go # 角色列表/创建/删除
│   │   └── channel/             # 频道服务器
│   │       ├── channel_server.go
│   │       └── handler/
│   │           ├── map.go       # 地图相关 (进入/移动/传送)
│   │           ├── chat.go      # 聊天
│   │           ├── combat.go    # 战斗
│   │           ├── inventory.go # 背包
│   │           ├── skill.go     # 技能
│   │           ├── quest.go     # 任务
│   │           ├── npc.go       # NPC 交互
│   │           ├── trade.go     # 交易
│   │           ├── party.go     # 组队
│   │           ├── friend.go    # 好友
│   │           └── cashshop.go  # 商城
│   ├── script/                  # JS 脚本引擎桥接层
│   │   ├── engine.go            # goja 引擎封装
│   │   ├── bindings.go          # Go→JS 函数绑定 (cm.xxx, pi.xxx 等)
│   │   └── loader.go            # 脚本加载/热加载
│   ├── service/                 # 业务逻辑层
│   │   ├── character_service.go
│   │   ├── combat_service.go
│   │   ├── drop_service.go
│   │   └── ...
│   └── api/                     # REST API (给管理后台)
│       ├── router.go            # gin/chi 路由
│       └── handler/
│           ├── account_api.go
│           └── character_api.go
├── scripts/                     # JS 脚本文件 (从原项目复制)
│   ├── npc/
│   ├── quest/
│   ├── event/
│   └── reactor/
├── config/
│   ├── config.go                # 配置结构体定义 + YAML 加载
│   └── config.yaml              # 默认配置文件
├── cli/                         # CLI 入口 (urfave/cli/v3)
│   ├── root.go                  # 根命令 + 全局 flag 定义
│   ├── serve.go                 # serve 子命令 (启动服务)
│   └── migrate.go               # migrate 子命令 (数据库迁移)
├── migrations/                  # SQL 迁移脚本 (从原项目提取)
│   └── 001_init.sql
├── web/                         # Vue.js 管理后台 (从 gms-ui/ 复制)
├── Makefile
├── go.mod
├── go.sum
└── README.md                    # ← 本文件
```

---

## 5. 五大核心问题及解决方案

### 5.1 协议层 — TCP 通信与封包

**问题**：冒险岛客户端使用自定义 TCP 协议，封包格式为 `[header: 4 bytes] [opcode: 2 bytes] [data: n bytes]`，服务端需要正确解析和构造封包。

**Go 方案**：

```go
// 一个封包的基本结构
type Packet struct {
    Opcode uint16
    Data   []byte
}

// 读取时：先读 4 字节包长 → 读完整包体 → 解密 → 解析 opcode
// 发送时：构造 body → 加密 → 计算长度 → 写入 TCP
```

**核心思路**：
- Go 标准库 `net` + `bufio` 即可，不需要 Netty
- 每个客户端连接用一个 goroutine 处理（goroutine 比 Java 线程轻量得多）
- 封包粘包/拆包用固定 4 字节包头（小端序，值 = 整个包的长度）

**关键代码路径**（参考原 Java 项目）：
- `net/server/` → `internal/network/`
- `tools/data/` → `internal/network/codec/`

---

### 5.2 加解密层

**问题**：GMS v0.83 使用自定义加密方案，封包在网络上传输时是加密的。

**加密流程**（推测，需对照原项目确认）：
```
客户端发送: 原始封包 → 加密 → TCP
服务端接收: TCP → 解密 → 原始封包 → 处理
```

**涉及算法**：
1. **AES** — Go 标准库 `crypto/aes` 直接支持
2. **Shanda 加密** — 冒险岛早期版本的自定义加密，本质是一种字节变换，约 30 行代码

**Go 实现要点**：
```go
// AES 加密器
type MapleCrypto struct {
    sendCipher    cipher.Stream    // AES/OFB 加密流
    recvCipher    cipher.Stream    // AES/OFB 解密流
    sendIV        []byte           // 发送端 IV (从客户端握手获取)
    recvIV        []byte           // 接收端 IV
    mapleVersion  uint16           // 版本号 (83)
}

// Shanda 加解密（可能在 v0.83 中已经淘汰，需确认）
func ShandaEncrypt(data []byte) []byte { ... }
func ShandaDecrypt(data []byte) []byte { ... }
```

**关键点**：加密初始化需要客户端握手阶段的 IV（通过 `getHello` 封包协商）。

---

### 5.3 游戏逻辑

**问题**：这是最大的一块。冒险岛有海量的游戏系统，每个都是一套逻辑。

**核心系统列表**（按实现优先级）：

| 优先级 | 系统 | 说明 |
|---|---|---|
| P0 | 登录认证 | 账号密码验证、PIN/PIC |
| P0 | 角色管理 | 创建、删除、选择角色、进入游戏 |
| P0 | 地图系统 | 进入地图、地图内移动、传送门 |
| P0 | 聊天 | 普通聊天、私聊、喇叭 |
| P1 | 背包系统 | 物品增删、装备穿戴、消耗品使用 |
| P1 | NPC 交互 | NPC 对话、商店 |
| P1 | 技能系统 | 技能学习、释放 |
| P1 | 战斗系统 | 攻击怪物、受到伤害、死亡 |
| P1 | 掉落系统 | 物品/金币掉落和拾取 |
| P1 | 任务系统 | 任务接取、进度、完成 |
| P2 | 组队 | 创建队伍、邀请、解散 |
| P2 | 好友 | 添加/删除好友、在线状态 |
| P2 | 交易 | 玩家间交易 |
| P2 | 商城 | Cash Shop |
| P3 | 公会 | 创建、管理、公会技能 |
| P3 | 活动 | 各种限时活动 |

**迁移策略**：
- **不要从零实现**。每个系统去原项目读 Java 代码，理解逻辑，翻译成 Go
- 优先实现 P0，让角色能登录、进地图、走动、聊天 → 这就是一个"可见的里程碑"
- 然后逐个啃 P1、P2、P3

**数据来源**：
- 游戏配置数据来自 WZ 文件（客户端资源），通过工具解析为 XML/JSON → Go 解析加载进内存
- 原项目有 `wz-zh-CN/` 等目录存放已解析的数据，可以直接复用

---

### 5.4 JS 脚本引擎

**问题**：原项目 61.5% 的代码是 JavaScript（NPC 脚本、任务脚本、事件脚本），在 Java 里通过 Nashorn 引擎执行。Go 需要一个 JS 引擎来执行这些脚本。

**✅ 已决策 — 方案：goja（暂时保留 JS）**

[goja](https://github.com/dop251/goja) 是纯 Go 实现的 ECMAScript 5.1 引擎，**无需 CGO**，性能足够。

**决策理由**：
- 桥接层工作量与替代方案（Lua）相同
- 可以零成本复用原项目全部现有 JS 脚本
- goja 成熟稳定，社区活跃

**未来备选 — Lua**：项目功能开发完成后，如果时间充裕，可以考虑将脚本迁移到 Lua（[gopher-lua](https://github.com/yuin/gopher-lua)）。Lua 是游戏行业脚本的事实标准，语法更简洁。JS → Lua 的状态机脚本可以写工具自动转换。但目前优先跑通功能，暂不折腾。

**核心挑战**：脚本中有大量调用 Java 桥接方法的调用，例如：

```javascript
// NPC 脚本示例 (Java 版)
var status = 0;
function start() {
    cm.sendNext("你好，冒险家！#b\r\n\r\n你想传送到哪里？");
}
function action(mode, type, selection) {
    status++;
    if (status == 1) {
        cm.sendSimple("#L0#前往金银岛#l\r\n#L1#前往天空之城#l");
    } else if (status == 2) {
        if (selection == 0) cm.warp(100000000, 0);
        else cm.warp(200000000, 0);
        cm.dispose();
    }
}
```

**桥接层方案**（`internal/script/bindings.go`）：

需要在 Go 中实现一个 `ConversationManager`（cm）和 `PlayerInteraction`（pi）等对象，把这些方法绑定到 goja 运行时：

```go
// Go 侧实现 cm 对象
func (e *Engine) setupBindings(vm *goja.Runtime, c *Client) {
    // 创建 cm 对象
    cm := vm.NewObject()
    cm.Set("sendNext", func(call goja.FunctionCall) goja.Value {
        msg := call.Argument(0).String()
        c.SendNPCDialog(NPCDialogNext, msg)
        return goja.Undefined()
    })
    cm.Set("sendSimple", func(call goja.FunctionCall) goja.Value { ... })
    cm.Set("warp", func(call goja.FunctionCall) goja.Value {
        mapId := call.Argument(0).ToInteger()
        portal := call.Argument(1).ToInteger()
        c.Warp(int(mapId), int(portal))
        return goja.Undefined()
    })
    cm.Set("dispose", func(call goja.FunctionCall) goja.Value { ... })
    // ... 更多方法

    vm.Set("cm", cm)
}
```

**需要实现的主要桥接对象**：
- `cm` (ConversationManager) — NPC 对话核心，约 100+ 方法
- `pi` (PlayerInteraction) — 玩家交互（组队任务等）
- `ms` (MapScriptManager) — 地图脚本
- `em` (EventManager) — 事件管理器

**迁移步骤**：
1. 先实现 `cm` 中最常用的 10-15 个方法（sendNext, sendSimple, warp, dispose, gainItem 等）
2. 跑通一个简单 NPC 脚本
3. 逐步补全其他方法

---

### 5.5 数据库设计

**✅ 已决策 — ORM：gorm**

**问题**：原项目使用 MySQL 8+，有一套完整的表结构。

**Go 方案**：
- **ORM**：`gorm` — Go 生态最流行的 ORM，类似 Hibernate 但更轻量
- **表结构**：从原项目的 SQL 初始化脚本提取，gorm AutoMigrate 或手写迁移 SQL
- **连接管理**：连接池由 gorm 内置管理

---

## 6. 分期实施计划

### 第一期：骨架搭建（预计 1-2 周）

**目标**：项目能编译，TCP 端口能监听，封包能正确加解密

- [x] 初始化 Go module 和项目目录结构
- [ ] 实现封包编解码 (`internal/network/codec/`)
- [ ] 实现 AES + Shanda 加解密 (`internal/crypto/`)
- [ ] 实现 TCP Server 和 Session 管理 (`internal/network/`)
- [ ] 实现 opcode 常量定义
- [ ] 编写 Makefile（`make run`, `make build`）
- [ ] 添加配置文件支持（urfave/cli + yaml.v3）

**里程碑**：客户端能建立连接，握手成功（能看到服务器在线）

---

### 第二期：登录与角色选择（预计 1-2 周）

**目标**：能登录、创建/选择角色、进入游戏

- [ ] 数据库连接和初始化
- [ ] Account/Character 模型和数据访问层
- [ ] Login Server 实现
  - [ ] 账号密码认证
  - [ ] 世界列表、频道列表
  - [ ] 角色列表、创建、删除
  - [ ] PIC/PIN 验证
- [ ] Channel Server 启动
  - [ ] 角色进入游戏流程

**里程碑**：客户端输入账号密码 → 看到角色列表 → 点角色 → 进入游戏地图

---

### 第三期：基础游戏体验（预计 3-4 周）

**目标**：角色能在地图中移动、聊天、打怪

- [ ] 地图系统
  - [ ] 加载地图数据 (WZ 解析后的数据)
  - [ ] 进入地图（角色刷新在他人视野中）
  - [ ] 地图内移动（走路、跳跃）
  - [ ] 传送门
  - [ ] 地图切换
- [ ] 聊天系统（普通/私聊/喇叭）
- [ ] 背包系统（查看、整理）
- [ ] NPC 交互
  - [ ] goja 引擎集成
  - [ ] cm 桥接层（核心 15 个方法）
  - [ ] 跑通第一个 NPC 对话脚本
- [ ] 战斗系统
  - [ ] 怪物刷新和 AI
  - [ ] 普通攻击
  - [ ] 技能释放
  - [ ] 伤害计算
  - [ ] 怪物死亡和掉落

**里程碑**：进去后能看到其他玩家、走路、聊天、打几只怪、跟 NPC 对话

---

### 第四期：核心玩法系统（预计 4-6 周）

**目标**：技能、任务、掉落等核心系统可用

- [ ] 技能系统完整实现
- [ ] 任务系统 (JS 脚本驱动)
- [ ] 掉落系统完整实现
- [ ] 组队系统
- [ ] 好友系统
- [ ] 交易系统

---

### 第五期：商城与社交（预计 2-3 周）

- [ ] Cash Shop
- [ ] 公会系统
- [ ] 结婚/家族等社交系统

---

### 第六期：管理后台 API（预计 1-2 周）

- [ ] REST API 实现（gin）
- [ ] Vue 管理后台对接

---

### 第七期：完善与优化

- [ ] 所有 JS 脚本兼容
- [ ] 性能优化
- [ ] 压力测试
- [ ] Docker 部署支持

---

## 7. 迁移策略

### 整体原则

1. **逐系统迁移，不贪大求全**：一个系统一个系统地对齐，每个系统写完要能跑
2. **先跑通骨架，再填肉**：第一期先让 TCP 通、加密对、客户端能握手
3. **Java 代码是参考书，不是规范**：理解逻辑后用 Go 惯用写法重写，不要逐行翻译
4. **JS 脚本尽量不动**：脚本保持原样，只改 Go 侧桥接层

### 每个系统的迁移步骤

```
1. 读 Java 源码 → 理解该系统的数据流和逻辑
2. 定义 Go 侧数据模型 (model)
3. 实现数据访问层 (store/DAO)
4. 实现业务逻辑 (service)
5. 实现封包处理 (handler) → 对接客户端
6. 对接 JS 脚本 (如果涉及)
7. 测试 → 客户端验证
```

### 封包对照法（关键技巧）

封包是协议层的"真理"。迁移每个系统时，在 Java 侧加日志打印收发的封包 hex dump，在 Go 侧也打印，两相对照，不一致就调 Go 侧。这样能快速定位问题。

---

## 8. 开发环境与工具

### 必备环境

| 工具 | 版本 | 用途 |
|---|---|---|
| Go | 1.21+ | 主力开发语言 |
| MySQL | 8.0+ | 数据库 |
| MapleStory 客户端 | GMS v0.83 | 测试用 |
| Git | latest | 版本控制 |

### 推荐工具

| 工具 | 用途 |
|---|---|
| Wireshark | 抓包分析协议（开发初期非常有用） |
| DBeaver | 数据库管理 |
| WZ 解析工具 | 提取客户端资源数据 |
| Postman | 调试 REST API |

---

## 9. 当前进度

> 最后更新: 2026-06-10

### 第一期：骨架搭建 ✅ 基本完成

- [x] 初始化 Go module 和项目目录结构
- [x] 实现封包编解码 (`internal/network/codec/`)
- [x] 实现 AES + 自定义混淆加解密 (`internal/crypto/`)，有真实抓包数据验证的测试用例
- [x] 实现 TCP Server 和 Session 管理 (`internal/network/`)
- [x] 实现 opcode 常量定义 (`internal/opcode/`)
- [x] 编写 Makefile
- [x] 配置文件支持（urfave/cli/v3 + yaml.v3）
- [x] 数据库连接 + gorm AutoMigrate (`internal/store/`)
- [x] 数据模型: Account, Character, Item, Skill (`internal/model/`)

**里程碑**: 客户端能建立连接，SERVER_HELLO → CLIENT_HELLO 握手成功，加密通道建立 ✅

### 第二期：登录与角色选择 🔄 进行中

- [x] SERVER_HELLO 握手（服务端先发言，明文交换 IV）
- [x] CLIENT_HELLO 解密 + IV 同步
- [x] 账号密码认证（bcrypt / SHA-1 / SHA-512 三种 hash 格式兼容）
- [x] 自动注册（`auto_register` 配置项）
- [x] 封禁检测 / 旧 hash 自动迁移到 bcrypt / 多开检测
- [x] LoginStatus 成功/失败响应封包构造
- [ ] 服务器列表响应 (HandleServerList)
- [ ] 角色列表 (HandleCharList)
- [ ] 角色创建 (HandleCharCreate)
- [ ] 角色选择进入游戏 (HandleCharSelect → Channel Server)
- [ ] 频道服务器握手（Channel 侧 SERVER_HELLO 流程）

**已知 Bug**: `SendPacket()` 中 `s.crypto.Encrypt(body)` 返回值被丢弃，响应包实际以明文发送，需要修

### 第三～七期：尚未开始

---

## 10. 待讨论问题

> 以下问题标记为"待定"，需要进一步讨论或研究后确定。

### 10.1 待确认的技术细节

- [x] **GMS v0.83 的加密方案**（2026-06-10 确认）：AES-256-OFB 流密码 + 外层 6 轮自定义混淆（EncryptData/DecryptData），见 `internal/crypto/maple_crypto.go`。Shanda 加密在该版本已淘汰，不再使用。
- [x] **封包格式**（2026-06-10 确认）：4 字节包头用 IV 编码（含版本号校验位和 XOR 混淆后的包体长度），不是简单的大端/小端长度。包头编码逻辑见 `EncodePacketHeader()`，解码见 `DecodePacketLength()`。
- [x] **握手流程**（2026-06-10 确认）：服务端先发言，发 16 字节明文 SERVER_HELLO（含 sendIV + recvIV）→ 客户端回 CLIENT_HELLO（加密，6 字节 = 4B header + 2B body）→ 后续所有通信加密。详见 `internal/crypto/crypto.go` 的 `GenServerHello()`。
- [ ] **WZ 数据解析**：现有的 wz-zh-CN/ 目录数据格式是什么（XML? JSON?），Go 如何加载

### 10.2 可以后续再决定的

- [ ] **配置热加载**：是否需要热加载，初期不必
- [ ] **脚本热加载**：NPC 脚本修改后是否需要重载（开发期很有用）

### 10.3 已决策事项

- [x] **脚本语言**：暂时保留 JS（goja），零成本复用现有脚本。未来可选迁移到 Lua（gopher-lua）
- [x] **单机 vs 拆分**：先单进程（Login + Channel 同进程），后续按需拆分
- [x] **ORM**：gorm（详见 [3.3 节](#33-数据库-orm--gorm)）
- [x] **REST API 框架**：gin（详见 [3.4 节](#34-rest-api-框架--gin)）
- [x] **CLI + 配置**：urfave/cli/v3 + gopkg.in/yaml.v3 直接解析（详见 [3.1](#31-cli-框架--urfavecliv3) + [3.2](#32-配置文件--yamlv3-直接解析) 节）
- [x] **日志**：logrus（详见 [3.5 节](#35-日志--logrus)）

### 10.4 CLI + 配置 使用示例

`urfave/cli/v3` 负责命令行参数，`yaml.v3` 负责解析配置文件，职责分离：

```go
// cli/root.go — 定义 flag
app := &cli.Command{
    Flags: []cli.Flag{
        &cli.StringFlag{Name: "config", Value: "config.yaml", Usage: "配置文件路径"},
    },
    Action: func(ctx context.Context, cmd *cli.Command) error {
        return runServer(cmd.String("config"))
    },
}

// config/config.go — 加载 YAML
func Load(path string) (*Config, error) {
    data, _ := os.ReadFile(path)
    var cfg Config
    yaml.Unmarshal(data, &cfg)
    return &cfg, nil
}
```
---

## 11. 注意事项
    1. 各种设计文档都存在当前目录的docs文件夹下。


---

## 附录：参考资料

- [BeiDouMS/BeiDou-Server](https://github.com/BeiDouMS/BeiDou-Server) — 原始 Java 项目
- [P0nk/Cosmic](https://github.com/P0nk/Cosmic) — BeiDou 的上游项目
- [urfave/cli/v3](https://github.com/urfave/cli) — CLI 框架
- [gopkg.in/yaml.v3](https://gopkg.in/yaml.v3) — YAML 解析
- [gorm](https://gorm.io/) — Go ORM
- [gin](https://github.com/gin-gonic/gin) — REST API 框架
- [logrus](https://github.com/sirupsen/logrus) — 日志库
- [goja](https://github.com/dop251/goja) — Go 的 JavaScript 引擎
- [gopher-lua](https://github.com/yuin/gopher-lua) — Go 的 Lua 引擎（未来备选）
- [MapleStory 协议文档](https://github.com/angelsl/ms-rebirth/wiki) — 社区协议参考（注意版本差异）

---

*最后更新: 2026-06-10*
