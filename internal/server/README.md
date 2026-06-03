# internal/server/

游戏服务器实例层，包含 Login Server 和 Channel Server 两个核心服务。

## 目录

| 目录 | 职责 | 端口 |
|---|---|---|
| `login/` | 登录服务器：账号认证、世界/频道列表、角色 CRUD、进入游戏 | 8484 (默认) |
| `channel/` | 频道服务器：地图、聊天、战斗、NPC、背包、任务、社交等所有游戏内逻辑 | 7575 (默认) |

## 启动流程

```
cli/serve.go
  → TCPServer.Listen(8484, login.HandleConnection)   // goroutine
  → TCPServer.Listen(7575, channel.HandleConnection)  // goroutine
  → 等待 SIGINT/SIGTERM
  → TCPServer.Shutdown()
```

两个服务在同一个进程、同一个 `TCPServer` 实例上运行，通过不同端口区分。后续可拆分为独立进程。
