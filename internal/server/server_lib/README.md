# internal/server/server_lib/

跨 login/channel 共用的工具库。封包构造、密码处理、世界/频道配置数据均集中在此。

**为什么放在 `server_lib/` 而非 `handler/`**：login 和 channel 两个服务都需要这些工具函数（如 `AfterLoginError` 既在 `AuthHandler.HandleCharSelect` 中用到，也在 `ChannelHandlerImpl.HandlePlayerLoggedin` 中用到）。放在独立的 lib 包中避免循环依赖。

## 文件

| 文件 | 职责 | 状态 |
|---|---|---|
| `packet.go` | 封包构造：`LoginStatusSuccess/Failed`、`ServerListEntry`、`CharList`、`CharNameResponse`、`AddNewCharEntry`、`ServerIP`、`AfterLoginError`、`GetCharInfo` (SET_FIELD) 等 | ✅ 已实现 |
| `password.go` | 密码处理：`HashPassword` (bcrypt)、`VerifyPassword`、`NeedsRehash` (SHA-1/SHA-512→bcrypt 迁移) | ✅ 已实现 |
| `world.go` | 世界/频道数据：`WorldInfo`、`ChannelInfo`、`WorldDataProvider` (世界列表、频道查找)、`RecommendedWorld` | ✅ 已实现 |

## GetCharInfo (SET_FIELD) — 0x7D

频道服务器最关键的一个封包。客户端收到它才会显示游戏画面。结构对齐 Java `PacketCreator.getCharInfo()`:

```
SET_FIELD:
  [channel-1:int] [battle:1] [updated:1] [reserved:2]
  [random1..3:int×3]
  addCharacterInfo:
    [unknown:long(-1)] [unknown:byte(0)]
    addCharStats: [id][name(13B)][gender][skin][face][hair][pets×3]...[spawnpoint]
    [buddyCapacity][linkedName(0)][meso]
    addInventoryInfo: [slots×5][time][equipped(空)][equip(空)][use(空)][setup(空)][etc(空)][cash(空)]
    addSkillInfo: [0][short(0)][short(0)]
    addQuestInfo: [short(0)][short(0)]
    addMiniGameInfo: [short(0)]
    addRingInfo: [short(0)][short(0)][short(0)]
    addTeleportInfo: [int(0)×15]
    addMonsterBookInfo: [int(0)][byte(0)][short(0)]
    addNewYearInfo: [short(0)]
    addAreaInfo: [short(0)]
    [terminator:short(0)]
  [timestamp:long]
```

## 未来规划

- [ ] `packet.go` 按功能域拆分（当前已近 500 行）：`login_packets.go`、`channel_packets.go`
- [ ] `GetCharInfo` 中背包/技能/任务数据从空壳变为真实数据填充
- [ ] 世界/频道配置从硬编码改为 YAML 配置文件驱动
- [ ] `packet.go` 添加更多封包：好友列表同步、家族/公会状态、Buff 恢复等
