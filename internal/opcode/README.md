# internal/opcode/

封包操作码（opcode）常量定义。

## 文件

| 文件 | 职责 |
|---|---|
| `login.go` | 登录服务器 opcode — 认证、服务器列表、角色 CRUD、状态码 |
| `channel.go` | 频道服务器 opcode — 地图、聊天、战斗、背包、NPC、任务、组队、好友、公会、商城 |

## 命名约定

```
LoginXxx       = 登录服务器相关
ChannelXxx     = 频道服务器相关
```

每个 opcode 注释标注方向：`接收 (Client → Server)` 或 `发送 (Server → Client)`。

## 注意

当前 opcode 值为**占位值**，需对照原 Java 项目逐一确认实际值。不同冒险岛版本的 opcode 不同（GMS v0.83 vs EMS vs MSEA 等）。
