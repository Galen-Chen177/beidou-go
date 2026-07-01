# internal/server/

游戏服务器实例层，包含 Login Server 和 Channel Server 两个核心服务。

## 目录

| 目录 | 职责 | 端口 | README |
|---|---|---|---|
| `login/` | 登录服务器：账号认证、世界/频道列表、角色 CRUD、进入游戏 | 8484 (默认) | [README](login/README.md) |
| `channel/` | 频道服务器：0x14 PLAYER_LOGGEDIN 进入游戏 → 即将扩展地图/战斗等 | 7575 (默认) | [README](channel/README.md) |
| `handler/` | 封包处理器实现：`AuthHandler` (login 接口)、`ChannelHandlerImpl` (channel 接口) | — | [README](handler/README.md) |
| `server_lib/` | 跨服务工具库：封包构造 (packet.go)、密码处理 (password.go)、世界/频道数据 (world.go) | — | [README](server_lib/README.md) |

## 架构关系

```
cli/serve.go (依赖注入)
  │
  ├─→ login.Server ──── 持有 login.LoginHandler 接口
  │     │                  ↑ 实现者: handler.AuthHandler
  │     │
  │     └─→ dispatch(0x01/0x04/0x05/0x06/0x13/0x15/0x16)
  │
  ├─→ channel.Server ── 持有 channel.ChannelHandler 接口
  │     │                  ↑ 实现者: handler.ChannelHandlerImpl
  │     │
  │     └─→ dispatch(0x14/0x23)
  │
  └─→ server_lib        ← login 和 channel 共用: packet.go, password.go, world.go
```

**接口反置原则**：接口定义在 `login/` 和 `channel/` 包中（靠近使用方），实现放在 `handler/` 包中（靠近依赖方）。这样 `login_server.go` 不依赖 `handler` 包，避免循环依赖。

## 启动流程

```
cli/serve.go
  → TCPServer.Listen(8484, login.HandleConnection)   // goroutine
  → TCPServer.Listen(7575, channel.HandleConnection)  // goroutine
  → 等待 SIGINT/SIGTERM
  → TCPServer.Shutdown()
```

两个服务在同一个进程、同一个 `TCPServer` 实例上运行，通过不同端口区分。后续可拆分为独立进程。

## 当前进度

| 模块 | 状态 |
|---|---|
| Login Server — 认证、服务器列表、角色 CRUD、角色选择 | ✅ 完成 |
| Channel Server — 握手、0x14 PLAYER_LOGGEDIN、GetCharInfo (SET_FIELD) | ✅ 完成 |
| 地图系统、聊天、战斗、NPC、背包、技能、任务 | ⬜ 待实现 |
| 好友、组队、公会、交易、商城 | ⬜ 待实现 |
