package opcode

// 频道服务器 接收 (Client → Server)
const (
	// 地图
	ChannelPlayerMove  uint16 = 0x1C // 玩家移动
	ChannelChangeMap   uint16 = 0x22 // 切换地图（传送门/脚本传送）
	ChannelEnterMap    uint16 = 0x19 // 进入地图（角色加载完成）
	ChannelMapTransfer uint16 = 0x22 // 地图间传送

	// 聊天
	ChannelChat uint16 = 0x2E // 普通聊天

	// 战斗
	ChannelPlayerAttack uint16 = 0x25 // 玩家攻击
	ChannelTakeDamage   uint16 = 0x26 // 受到伤害
	ChannelSkillUse     uint16 = 0x2B // 使用技能

	// 背包
	ChannelInventoryOp uint16 = 0x47 // 背包操作
	ChannelChangeEquip uint16 = 0x29 // 更换装备

	// NPC
	ChannelNpcTalk     uint16 = 0x3A // NPC 对话
	ChannelNpcTalkMore uint16 = 0x3C // NPC 对话继续
	ChannelNpcShop     uint16 = 0x3B // NPC 商店

	// 任务
	ChannelQuestAction uint16 = 0x80 // 任务操作

	// 交互
	ChannelPlayerInteraction uint16 = 0x44 // 玩家交互（交易/组队邀请等）
	ChannelPartyOperation    uint16 = 0x6C // 组队操作
	ChannelFriendOperation   uint16 = 0x66 // 好友操作
	ChannelGuildOperation    uint16 = 0x51 // 公会操作

	// 其他
	ChannelHello         uint16 = 0x14 // 频道服务 cli hello
	ChannelChangeChannel uint16 = 0x18 // 切换频道
	ChannelCashShop      uint16 = 0x86 // 商城操作
	ChannelPlayerUpdate  uint16 = 0x40 // 玩家信息更新请求
)

// 频道服务器 发送 (Server → Client)
const (
	// 核心
	ChannelSetField uint16 = 0x7D // SET_FIELD (getCharInfo — 角色进入游戏的关键封包)

	// 地图
	ChannelWarpToMap      uint16 = 0x3E // 传送到地图
	ChannelSpawnPlayer    uint16 = 0x3D // 刷新其他玩家
	ChannelRemovePlayer   uint16 = 0x15 // 移除玩家
	ChannelPlayerMovement uint16 = 0x8F // 广播玩家移动
	ChannelSpawnNPC       uint16 = 0x7F // 刷新 NPC
	ChannelSpawnMonster   uint16 = 0x86 // 刷新怪物
	ChannelKillMonster    uint16 = 0x88 // 怪物死亡

	// 聊天
	ChannelChatMsg uint16 = 0x8B // 聊天消息

	// 状态
	ChannelPlayerHP    uint16 = 0x3F // 玩家 HP/MP 更新
	ChannelStatChanged uint16 = 0x1E // 属性变化
	ChannelSkillEffect uint16 = 0x5B // 技能特效

	// NPC
	ChannelNpcTalkResp     uint16 = 0x9D // NPC 对话响应
	ChannelNpcShopResp     uint16 = 0x7B // NPC 商店响应
	ChannelNpcActionResult uint16 = 0x83 // NPC 操作结果

	// 背包
	ChannelInventoryGrow   uint16 = 0x3A // 背包变化通知
	ChannelInventoryOpResp uint16 = 0x48 // 背包操作结果

	// 提示消息
	ChannelMessage uint16 = 0x20 // 顶部提示消息
	ChannelDropTip uint16 = 0x19 // 掉落提示

	// 服务器消息（系统通知）
	ChannelServerMessage uint16 = 0x44 // 服务器公告/喇叭

	// 键位 & 快捷栏
	ChannelKeymap        uint16 = 0x14F // KEYMAP
	ChannelQuickslotInit uint16 = 0x9F  // QUICKSLOT_INIT
	ChannelMacroSysData  uint16 = 0x7C  // MACRO_SYS_DATA_INIT

	// 自动喝药
	ChannelAutoHpPot uint16 = 0x150 // AUTO_HP_POT
	ChannelAutoMpPot uint16 = 0x151 // AUTO_MP_POT
)
