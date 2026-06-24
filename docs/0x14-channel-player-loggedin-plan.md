# 0x14 (PLAYER_LOGGEDIN) 频道服务器 — 实现计划

> 创建日期: 2026-06-24
> 状态: 待实现
> 对话摘要: 用户询问 Java 端 0x13 的处理，发现 0x13 是 Login Server 的 CHAR_SELECT，频道服务器对应的是 0x14 PLAYER_LOGGEDIN。当前 Go 项目 channel_server 的 dispatch 是从 login_server 复制过来的，需要改造。

---

## 0x13 vs 0x14 — 关键区别

| Opcode | 名称 | 服务端 | Java 文件 | 职责 |
|--------|------|--------|-----------|------|
| `0x13` | CHAR_SELECT | **Login Server** | `handlers/login/CharSelectedHandler.java` | 验证角色归属，发送频道 IP:Port (getServerIP) |
| `0x14` | PLAYER_LOGGEDIN | **Channel Server** | `channel/handlers/PlayerLoggedinHandler.java` | 加载角色数据，发送 getCharInfo，角色进入世界 |

### 数据流

```
Client                     Login Server (8484)              Channel Server (8485)
  │                              │                                │
  │── 0x13 CHAR_SELECT ────────→│                                │
  │                              │ 验证角色、查IP:Port             │
  │←── getServerIP (0x0C) ──────│                                │
  │                              │                                │
  │  (断开Login, 连接Channel)     │                                │
  │                                                               │
  │── 0x14 PLAYER_LOGGEDIN ────────────────────────────────────→│
  │                                                               │ 加载角色、入地图
  │←── SET_FIELD (0x7D) getCharInfo ──────────────────────────│
  │                                                               │
```

---

## 当前 Go 项目存在的问题

### 1. `channel_handler.go` 接口错误
[channel_handler.go](../internal/server/channel/channel_handler.go) 的 `ChannelHandler` 接口和 `LoginHandler` 完全一样，包含 `HandleCheckPassword`、`HandleServerList` 等登录方法。频道服务器不应该有这些方法。

### 2. `channel_server.go` dispatch 错误
[channel_server.go](../internal/server/channel/channel_server.go) 的 `dispatch()` 路由的是 login opcode（0x01, 0x04, 0x05, 0x13, 0x15, 0x16...），而 `ChannelHello` (0x14) 只有一个空注释 `// cli hello`。

### 3. `clic/serve.go` 共用 handler
[serve.go](../clic/serve.go) 将同一个 `authHandler` 传给了 login server 和 channel server。频道服务器需要自己独立的 handler 实现。

---

## 实现方案

### Step 1: 修改 `ChannelHandler` 接口

**文件:** `internal/server/channel/channel_handler.go`

```go
package channel

import "beidou-go/internal/network"

type ChannelHandler interface {
    HandlePlayerLoggedin(sess *network.Session, charID int32)
}
```

删除所有 login 相关方法。

### Step 2: 添加 SendOpcode 常量

**文件:** `internal/opcode/channel.go`

新增发送 opcode：

```go
const (
    // 频道服务器 发送 (Server → Client)
    ChannelSetField            uint16 = 0x7D   // SET_FIELD (getCharInfo)
    ChannelSelectCharByVac     uint16 = 0x09   // SELECT_CHARACTER_BY_VAC (进入游戏错误)
    ChannelServerMessage       uint16 = 0x44   // SERVERMESSAGE
    ChannelBuddylist           uint16 = 0x3F   // BUDDYLIST
    ChannelFamilyPrivilegeList uint16 = 0x64   // FAMILY_PRIVILEGE_LIST
    ChannelFamilyInfoResult    uint16 = 0x5F   // FAMILY_INFO_RESULT
    ChannelSetGender           uint16 = 0x3A   // SET_GENDER
    ChannelClaimStatusChanged  uint16 = 0x2F   // CLAIM_STATUS_CHANGED (enableReport)
    ChannelKeymap              uint16 = 0x14F  // KEYMAP
    ChannelQuickslotInit       uint16 = 0x9F   // QUICKSLOT_INIT
    ChannelMacroSysDataInit    uint16 = 0x7C   // MACRO_SYS_DATA_INIT
    ChannelAutoHpPot           uint16 = 0x150  // AUTO_HP_POT
    ChannelAutoMpPot           uint16 = 0x151  // AUTO_MP_POT
    ChannelUpdateHpMpAlert     uint16 = 0x1000 // UPDATE_HPMPAALERT
    ChannelGuildOperation      uint16 = 0x7E   // GUILD_OPERATION
    ChannelAllianceOperation   uint16 = 0x8F   // ALLIANCE_OPERATION
    ChannelParcel              uint16 = 0x51   // PARCEL (duey)
    ChannelScriptProgressMsg   uint16 = 0x8A   // SCRIPT_PROGRESS_MESSAGE
)
```

### Step 3: 构建 `GetCharInfo` 封包

**文件:** `internal/server/server_lib/packet.go`

这是最关键的封包——客户端收到它才会显示游戏画面。

对齐 Java `PacketCreator.getCharInfo()`:

```
格式 (不含 opcode, SetField = 0x7D):
  [channel-1: int]        — 频道号-1 (客户端用0-based)
  [battle: byte(1)]       — 是否战斗中 (写死1=否)
  [updated: byte(1)]      — 是否更新 (写死1)
  [reserved: short(0)]    — 保留
  [random1: int]          — 随机种子1
  [random2: int]          — 随机种子2
  [random3: int]          — 随机种子3
  ── addCharacterInfo ──
  [unknown: long(-1)]     — 未知标志
  [unknown: byte(0)]      — 未知标志
  ── addCharStats ──
  [charID: int]
  [name: padded_string(13)]
  [gender: byte]
  [skinColor: byte]
  [face: int]
  [hair: int]
  [petID1: long(0)] [petID2: long(0)] [petID3: long(0)]
  [level: byte]
  [job: short]
  [str: short] [dex: short] [int: short] [luk: short]
  [hp: short] [maxhp: short] [mp: short] [maxmp: short]
  [remainingAp: short]
  [remainingSp: short]
  [exp: int]
  [fame: short]
  [gachaExp: int]
  [mapID: int]
  [spawnPoint: byte]
  [reserved: int(0)]
  ── 好友容量 ──
  [buddyCapacity: byte]
  ── 关联名 ──
  [linkedName: byte(0)]   — 无关联名
  ── 金币 ──
  [meso: int]
  ── addInventoryInfo (MVP: 空) ──
  [equip: 见下文]
  [use: 见下文]
  [setup: 见下文]
  [etc: 见下文]
  [cash: 见下文]
  ── addSkillInfo (MVP: 空) ──
  ── addQuestInfo (MVP: 空) ──
  ── addMiniGameInfo (全0) ──
  ── addRingInfo (全0) ──
  ── addTeleportInfo (全0) ──
  ── addMonsterBookInfo (全0) ──
  ── addNewYearInfo (全0) ──
  ── addAreaInfo (全0) ──
  [terminator: short(0)]
  ── 时间戳 ──
  [timestamp: long]       — 服务器当前毫秒时间戳
```

**空的 Inventory 格式:**
```
[equipSlots: byte] [useSlots: byte] [setupSlots: byte] [etcSlots: byte] [cashSlots: byte]
[equip: 每个装备 1B(类型) + 1B(slot) + 3*short + 7*int + string(name) + ...]
0xFF ← equip 结束标记
[use: 对每个非空格 1B(slot) + 数据]
0xFF ← use 结束标记
... 同理 setup, etc, cash
```

MVP 阶段所有背包都为空（只写 slot 数量 + 0xFF 结束标记）。

**空的 Skill 格式:**
```
[skillCount: short(0)]
```

**空的 Quest 格式:**
```
[ongoingCount: short(0)]
[completedCount: short(0)]
```

**其他所有 info 部分都只写 `short(0)` + `byte(0)` 结束。**

### Step 4: 重写 `channel_server.go` 的 dispatch

**文件:** `internal/server/channel/channel_server.go`

删除所有 login opcode 的 case，替换为：

```go
func (s *Server) dispatch(sess *network.Session, packet *codec.Packet) {
    switch packet.Opcode {
    case opcode.ChannelHello: // 0x14 PLAYER_LOGGEDIN
        if len(packet.Data) >= 4 {
            charID := int32(packet.Data[0]) | int32(packet.Data[1])<<8 |
                int32(packet.Data[2])<<16 | int32(packet.Data[3])<<24
            s.handler.HandlePlayerLoggedin(sess, charID)
        }
    case opcode.LoginClientHello: // 0x23 CLIENT_HELLO 握手
        s.log.Infof("[channel] client hello handshake complete")
    default:
        s.log.Warnf("[channel] 未处理的 opcode: session_id=%d, opcode=0x%04X", sess.ID(), packet.Opcode)
    }
}
```

### Step 5: 创建 `ChannelHandlerImpl`

**新文件:** `internal/server/handler/channel_handler.go`

```go
package handler

import (
    "github.com/sirupsen/logrus"
    "beidou-go/config"
    "beidou-go/internal/model"
    "beidou-go/internal/network"
    "beidou-go/internal/server/server_lib"
    "beidou-go/internal/store"
)

type ChannelHandlerImpl struct {
    characterStore *store.CharacterStore
    cfg            *config.Config
    log            *logrus.Logger
}

func NewChannelHandler(
    characterStore *store.CharacterStore,
    cfg *config.Config,
    log *logrus.Logger,
) *ChannelHandlerImpl {
    return &ChannelHandlerImpl{
        characterStore: characterStore,
        cfg:            cfg,
        log:            log,
    }
}

func (h *ChannelHandlerImpl) HandlePlayerLoggedin(sess *network.Session, charID int32) {
    h.log.Infof("[Channel] HandlePlayerLoggedin: session=%d, charID=%d", sess.ID(), charID)

    // 1. 从数据库加载角色
    chr, err := h.characterStore.FindByID(charID)
    if err != nil {
        h.log.Warnf("[Channel] 角色不存在: charID=%d, err=%v", charID, err)
        sess.SendPacket(server_lib.GetAfterLoginError(17))
        return
    }

    // 2. 设置 session 关联数据
    sess.CharID = charID
    sess.AccountID = uint(chr.Accountid)
    sess.WorldID = byte(chr.World)

    // 3. 发送 GetCharInfo (SET_FIELD) — 这是最关键的包
    channelID := byte(1)
    if sess.ChannelID > 0 {
        channelID = sess.ChannelID
    }
    if err := sess.SendPacket(server_lib.GetCharInfo(chr, channelID)); err != nil {
        h.log.Errorf("[Channel] 发送 GetCharInfo 失败: %v", err)
        return
    }

    h.log.Infof("[Channel] 角色进入游戏成功: name=%s, id=%d, map=%d", chr.Name, chr.ID, chr.Map)
}
```

### Step 6: 修改 `clic/serve.go` 的注入

将：
```go
channelSrv := channel.NewServer(cfg, tcpSrv, log, authHandler)
```

改为：
```go
channelHandler := handler.NewChannelHandler(characterStore, cfg, log)
channelSrv := channel.NewServer(cfg, tcpSrv, log, channelHandler)
```

---

## Java 参考代码路径

| Java 文件 | 作用 |
|-----------|------|
| `net/opcodes/RecvOpcode.java:42` | `CHAR_SELECT(0x13)`, `PLAYER_LOGGEDIN(0x14)` |
| `net/opcodes/SendOpcode.java:154` | `SET_FIELD(0x7D)` |
| `net/PacketProcessor.java:148` | Login 注册 `CharSelectedHandler` 处理 0x13 |
| `net/PacketProcessor.java:183` | Channel 注册 `PlayerLoggedinHandler` 处理 0x14 |
| `handlers/login/CharSelectedHandler.java` | 0x13 处理：验证→发 getServerIP |
| `channel/handlers/PlayerLoggedinHandler.java` | 0x14 处理：加载角色→发 getCharInfo→广播 |
| `util/PacketCreator.java:957` | `getCharInfo()` — SET_FIELD 封包 |
| `util/PacketCreator.java:223` | `addCharacterInfo()` — 全部角色数据 |
| `util/PacketCreator.java:172` | `addCharStats()` — 角色属性 |

---

## 验证步骤

1. `go build ./...` — 必须编译通过
2. `go run . serve` — 启动服务器
3. 用 MapleStory v0.83 客户端：
   - 登录账号 → 看到角色列表
   - 点击角色 → 客户端应连接到频道服务器并进入游戏画面
4. 检查服务端日志：应看到 `[channel] opcode=0x0014` 和 `[Channel] 角色进入游戏成功`

---

## 后续迭代 (不在本次范围)

- [ ] 好友列表同步
- [ ] 家族/公会/组队状态恢复
- [ ] 宠物/坐骑恢复
- [ ] Buff/Debuff/技能冷却恢复
- [ ] 背包/技能数据填充到 getCharInfo
- [ ] GM 自动隐身
- [ ] Duey 包裹通知
- [ ] NPC Scriptable 配置
- [ ] 事件恢复 (EventRecall)
- [ ] 结婚伴侣上线通知
