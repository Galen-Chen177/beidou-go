# internal/server/handler/

封包处理器实现层。这里的结构体实现了 `login.ChannelHandler` / `channel.ChannelHandler` 等接口，被 `login.Server` / `channel.Server` 通过依赖注入引用。

接口定义在 `login/` 和 `channel/` 包中（而非 handler 包），目的是**避免 import cycle**：`login_server.go` 持有接口 → 不依赖 handler 包；`handler/*.go` 实现接口 → 可以 import login/channel 包。

## 文件

| 文件 | 职责 | 状态 |
|---|---|---|
| `auth.go` | `AuthHandler` — 登录认证处理器：密码验证、服务器列表、角色 CRUD、角色选择。实现 `login.LoginHandler` 接口 | ✅ 已实现 |
| `channel_handler.go` | `ChannelHandlerImpl` — 频道处理器：0x14 PLAYER_LOGGEDIN 角色进入游戏。实现 `channel.ChannelHandler` 接口 | ✅ 已实现 |

## 设计原则

- **每个服务对应一个 HandlerImpl**：Login 用 `AuthHandler`，Channel 用 `ChannelHandlerImpl`。不共用。
- **依赖最小化**：HandlerImpl 只依赖 `store.*Store`（数据访问）、`config.Config`（配置）、`logrus.Logger`（日志），不引入网络层细节。
- **SendPacket 在 handler 内部完成**：Handler 调用 `server_lib.*()` 构造封包 → `sess.SendPacket()` 发送。封包构造逻辑集中在 `server_lib/` 包中。

## 未来规划

- [ ] `channel_handler.go` 按功能域拆分：`map.go`、`chat.go`、`combat.go`、`inventory.go`、`skill.go`、`quest.go`、`npc.go` 等
- [ ] 每个子 handler 通过组合方式注入 `ChannelHandlerImpl`，避免单个文件膨胀
- [ ] `auth.go` 可能拆出 `world.go`（世界/频道逻辑）和 `character.go`（角色创建逻辑）
