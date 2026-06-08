# 登陆处理 — Go 侧实现计划

## Context

当前项目已搭建好骨架（TCP server、封包编解码、AES/OFB 加密、Session 管理、握手流程），需要实现第二阶段核心功能：**登陆认证**，让客户端能成功输入账号密码 → 看到角色列表。

参考：Java 源码分析文档 [登陆验证逻辑分析.md](登陆验证逻辑分析.md)，已梳理清楚完整链路。

## 实现文件清单（按依赖顺序）

### 1. `internal/crypto/aes.go` — 修复 IV→AES Key 扩展（bugfix）

**问题**：当前 `NewMapleCrypto` 直接把 4 字节 IV 传给 `aes.NewCipher()`，但 AES-128 要求 16 字节密钥，会导致 panic。

**修改**：新增 `expandIV(iv []byte) []byte`，将 4 字节 IV 重复 4 次得到 16 字节密钥。修改 `NewMapleCrypto` 在创建 cipher 前调用 expansion。

**原因**：这是登陆流程的前置依赖——握手完成后 ReadPacket/SendPacket 需要加解密。

### 2. `internal/model/account.go` — 补全字段

当前缺少 Java `accounts` 表的关键字段。新增：

| 字段 | Go 类型 | DB column |
|------|---------|-----------|
| Gender | int8 | gender |
| CharSlots | int8 | characterslots |
| TOS | bool | tos |
| Language | string | language |

### 3. `internal/store/account_store.go` — 新建 Account DAO

提供 3 个方法：

- `FindAccountByName(name string) (*model.Account, error)` — `WHERE name = ?`
- `CreateAccount(account *model.Account) error` — 自动注册用
- `UpdatePassword(id int32, hash string) error` — bcrypt 迁移用

### 4. `internal/server/login/password.go` — 新建密码校验

- `VerifyPassword(plain, hash string) bool` — 判断 hash 前缀（`$2a$`/`$2b$`→bcrypt，40 字符 hex→SHA-1，128 字符 hex→SHA-512），调用对应算法
- `HashPassword(plain string) (string, error)` — bcrypt 生成（cost=10）
- `NeedsRehash(hash string) bool` — 非 bcrypt 格式返回 true

依赖：`golang.org/x/crypto/bcrypt`（已在 go.sum 中）。

### 5. `internal/server/login/session_coordinator.go` — 新建多开检测

简单的内存级实现：

- `type SessionCoordinator struct` 内部用 `sync.Map`（key=账号名, value=sessionID）
- `AttemptLogin(name string, sessionID uint32) error` — 已存在则返回错误
- `Logout(name string)` — 移除条目

### 6. `internal/server/login/packet.go` — 新建封包构造器

提供构造 LoginStatus 响应封包的函数：

- `LoginStatusSuccess(account *model.Account) *codec.Packet` — opcode=0x00, status=0, 后接 account_id/gender/gm/name
- `LoginStatusFailed(reason byte) *codec.Packet` — opcode=0x00, status=error_code

封包格式（参考 GMS v0.83 社区文档）：
```
成功: [opcode 0x00][status:0][account_id:4B LE][gender:1B][gm:1B][name:str]
失败: [opcode 0x00][status:error_code][reason:0]
```

### 7. `internal/server/login/login_server.go` — 改造握手 + 路由分发

**修改点 1**：CLIENT_HELLO 解析
- 从 `readHandshakePacket` 返回的原始字节中解析出客户端的 sendIV/recvIV（格式与 SERVER_HELLO 一致：[2B 长度][2B 版本][2B 子版本][4B sendIV][4B recvIV]...）
- 调用 `crypto.NewMapleCrypto(sendIV, recvIV, 83)` 创建加解密器
- 调用 `sess.SetCrypto(mc)` 注入 Session

**修改点 2**：封包读取循环中添加 opcode 路由
- opcode `0x01` (LOGIN_PASSWORD) → `authHandler.HandleCheckPassword()`
- 其他 opcode 打印未处理日志

### 8. `internal/server/login/handler/auth.go` — 重写登陆处理

改造 `AuthHandler` 结构体，注入依赖：

```go
type AuthHandler struct {
    store       *store.AccountStore     // (或直接传 *gorm.DB)
    coordinator *SessionCoordinator
    log         *logrus.Logger
}
```

**`HandleCheckPassword(sess *network.Session, data []byte)` 实现**：

```
1. 用 codec.Reader 解析 data:
   - ReadString() → login (账号名)
   - ReadString() → password (密码明文)
   - Skip(6)     → 跳过填充
   - ReadBytes(4) → hwid

2. 查库: store.FindAccountByName(login)

3. 账号不存在:
   → 如开启 auto_register: 创建账号(bcrypt hash)，继续
   → 否则: Send LoginStatus(5 "未注册")，return

4. 封禁检测: if account.Banned → Send LoginStatus(3 "已封禁")，return

5. 密码校验: VerifyPassword(password, account.Password)
   → 失败: Send LoginStatus(4 "密码错误")，return

6. 旧 hash 迁移: if NeedsRehash → UpdatePassword(account.ID, bcrypt)

7. 多开检测: coordinator.AttemptLogin(login, sess.ID())
   → 失败: Send LoginStatus(7 "已在其他位置登录")，return

8. 成功:
   - sess.AccountID = account.ID
   - Send LoginStatusSuccess(account)
```

### 9. `cli/serve.go` — 启用数据库

取消注释 `store.InitDB(cfg.Database)` 调用，并将 `*gorm.DB` 实例传递给 login handler。

### 10. `config/config.go` — 新增 `auto_register` 配置项

在 `LoginConfig` 中加 `AutoRegister bool \`yaml:"auto_register"\``，配置文件对应加项（默认 true，便于初期测试）。

## 构建顺序

```
Step 1: crypto/aes.go (fix)        → 基础依赖
Step 2: model/account.go (补字段)   → 数据模型
Step 3: store/account_store.go     → 数据访问
Step 4: login/password.go          → 密码工具
Step 5: login/session_coordinator  → 多开检测
Step 6: login/packet.go            → 响应封包
Step 7: login/handler/auth.go      → 核心处理器
Step 8: login/login_server.go      → 改造握手+路由
Step 9: config + cli/serve.go      → 集成启动
```

## 验证方式

1. `go build ./...` — 确保编译通过
2. 启动 `beidou-go serve`，数据库就绪
3. 用冒险岛 GMS v0.83 客户端连接 `127.0.0.1:8484`：
   - 输入正确账号密码 → 看到 LoginStatus success（进入角色选择）
   - 输入错误密码 → 提示密码错误
   - 不存在的账号 → 提示未注册（或自动注册）
4. 观察服务端日志：每个封包的收发都有 hex dump（debug 级别）
