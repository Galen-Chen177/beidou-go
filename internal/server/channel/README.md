# internal/server/channel/

频道服务器，负责所有游戏内逻辑。这是整个项目最复杂的模块。

## 文件

| 文件 | 职责 | 状态 |
|---|---|---|
| `channel_server.go` | `Server` — 启动端口监听、握手 (SERVER_HELLO/CLIENT_HELLO)、加解密通道建立、dispatch() 按 opcode 分发 | ✅ 已实现 |
| `channel_handler.go` | `ChannelHandler` 接口 — 定义 `HandlePlayerLoggedin(sess, charID)` | ✅ 已实现 |

## dispatch 路由

```
0x23 CLIENT_HELLO → 握手完成，IV 同步
0x14 PLAYER_LOGGEDIN → HandlePlayerLoggedin(sess, charID)
default → Warn("[channel] 未处理的 opcode")
```

## 数据流

```
Client                     Login Server (8484)              Channel Server (8485)
  │                              │                                │
  │── 0x13 CHAR_SELECT ────────→│                                │
  │                              │ 验证角色、查 IP:Port             │
  │←── getServerIP (0x0C) ──────│                                │
  │                              │                                │
  │  (断开 Login，连接 Channel)    │                                │
  │── 0x14 PLAYER_LOGGEDIN ────────────────────────────────────→│
  │                                                               │ 查 DB 加载角色
  │                                                               │ 设 session 关联
  │←── SET_FIELD (0x7D) getCharInfo ───────────────────────────│
  │                                                               │
```

## handler/ (实现放在 internal/server/handler/)

接口方法的具体实现在 `internal/server/handler/channel_handler.go`（`ChannelHandlerImpl`），不在本目录。这是为了避免 import cycle：channel 包定义接口，handler 包实现接口。

## 未来规划

- [ ] 地图系统：进入地图 (`0x19`)、移动 (`0x1C`)、传送门 (`0x22`)、地图切换
- [ ] 聊天系统：普通聊天 (`0x2E`)、私聊、喇叭
- [ ] 战斗系统：普攻 (`0x25`)、技能释放 (`0x2B`)、伤害计算、怪物死亡
- [ ] 背包系统：物品增删 (`0x47`)、装备穿戴 (`0x29`)
- [ ] NPC 交互：对话 (`0x3A`)、商店 (`0x3B`)
- [ ] 技能系统：学习、加点
- [ ] 任务系统：接取、进度 (`0x80`)
- [ ] 社交系统：好友 (`0x66`)、组队 (`0x6C`)、公会 (`0x51`)、交易 (`0x44`)
- [ ] 商城：进入商城 (`0x86`)
- [ ] 频道切换 (`0x18`)
- [ ] `GetCharInfo` 中填充真实背包/技能/任务数据（非空壳）
- [ ] 好友列表同步、家族/公会状态恢复、Buff/Debuff 恢复
- [ ] GM 自动隐身
