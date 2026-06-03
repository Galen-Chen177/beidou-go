# internal/server/channel/

频道服务器，负责所有游戏内逻辑。这是整个项目最复杂的模块。

## 文件

| 文件 | 职责 |
|---|---|
| `channel_server.go` | `Server` — 频道服务器主逻辑：`Start()` 注册端口监听，`handleConnection()` 处理每个进入游戏的客户端（当前为骨架） |

## handler/ (待实现)

按功能域分子文件：

| 文件 | 职责 |
|---|---|
| `map.go` | 地图系统：进入地图、移动、传送门、地图切换 |
| `chat.go` | 聊天：普通、私聊、喇叭 |
| `combat.go` | 战斗：普攻、技能伤害计算、怪物死亡 |
| `inventory.go` | 背包：物品增删、装备穿戴 |
| `skill.go` | 技能：学习、加点、释放 |
| `quest.go` | 任务：接取、进度更新、完成 |
| `npc.go` | NPC 交互：对话、商店 |
| `trade.go` | 交易：玩家间物品/金币交易 |
| `party.go` | 组队：创建、邀请、解散 |
| `friend.go` | 好友：添加/删除、在线状态 |
| `cashshop.go` | 商城：进入商城、购买 |
