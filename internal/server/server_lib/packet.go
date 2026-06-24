package server_lib

import (
	"beidou-go/internal/model"
	"beidou-go/internal/network/codec"
	"beidou-go/internal/opcode"
)

// ──────────── LoginStatus (0x00) ────────────

// LoginStatusSuccess 构造登录成功响应封包 (包4)
//
//	对应 Java: PacketCreator.getAuthSuccess()
//
//	解密后数据结构 (共 47+ 字节，不含 opcode):
//	  [reserved: int(0)]          — 保留位，写死 0
//	  [reserved: short(0)]        — 保留位，写死 0
//	  [accID: int]                — 账号数据库 ID (LE)
//	  [gender: byte]              — 性别 (0=男, 1=女)
//	  [isGM: bool]                — 是否 GM 账号
//	  [adminByte: byte]           — 管理员标识 (GM=0x80, 普通=0x00)
//	  [country: byte]             — 国家代码 (写死 0)
//	  [accountName: string]       — 账号名 (2B 长度前缀 + ASCII)
//	  [terminator: byte(0)]       — 字符串结束符
//	  [quietBan: byte]            — 静默封禁标志 (0=否)
//	  [quietBanTime: long(0)]     — 静默封禁时间戳
//	  [createdAt: long]           — 账号创建时间戳 (毫秒, LE)
//	  [removePrompt: int(1)]      — 1=移除"选择你想进入的世界"提示
//	  [pinEnabled: byte]          — PIN 码状态: 0=需要设置, 1=已禁用
//	  [picEnabled: byte]          — PIC 状态: 0=需要注册, 1=需要输入, 2=已禁用
func LoginStatusSuccess(account *model.Account) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginStatus)

	// ── 保留位 ──
	w.WriteInt(0)   // reserved int
	w.WriteShort(0) // reserved short

	// ── 账号基本信息 ──
	w.WriteInt(uint32(account.ID))    // 账号 ID
	w.WriteByte(byte(account.Gender)) // 性别 (0=男, 1=女)

	// ── GM 标识 ──
	// 对齐 Java: (use_enforce_admin_account || canFly) && getGMLevel() > 1
	// 简化实现：webadmin > 1 视为 GM (webadmin=1 是见习GM，不发adminByte)
	isGM := account.Webadmin > 1
	w.WriteBool(isGM) // isGM flag
	if isGM {
		w.WriteByte(0x80) // admin byte (GM)
	} else {
		w.WriteByte(0x00) // admin byte (普通)
	}

	// ── 国家代码 ──
	w.WriteByte(0)

	// ── 账号名 + 结束符 ──
	w.WriteString(account.Name)
	w.WriteByte(0) // null terminator

	// ── 静默封禁 ──
	w.WriteByte(0) // quietBan = false
	w.WriteLong(0) // quietBan timestamp

	// ── 账号创建时间 ──
	// 对齐 Java: p.writeLong(0) 硬编码 0
	w.WriteLong(0)

	// ── 移除"选择你想进入的世界"提示 ──
	w.WriteInt(1)

	// ── PIN / PIC 状态 (暂无二次密码功能，均设为已禁用) ──
	w.WriteByte(1) // pinEnabled: 0=需要设置PIN, 1=PIN已禁用
	w.WriteByte(2) // picEnabled: 0=需要注册PIC, 1=需要输入PIC, 2=PIC已禁用

	return w.Packet()
}

// LoginStatusFailed 构造登录失败响应封包
//
//	格式: [status:error_code][reason:0]
func LoginStatusFailed(reason byte) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginStatus)

	w.WriteByte(reason) // status = error code
	w.WriteByte(0)      // reason (0 = no additional info)

	return w.Packet()
}

// ──────────── ServerList (0x0A) ────────────

// ServerListEntry 构造服务器列表封包 (包6)
//
//	对应 Java: PacketCreator.getServerList()
//
//	解密后数据结构:
//	  [serverID: byte]         — 服务器 ID
//	  [serverName: string]     — 服务器名称
//	  [flag: byte]             — 服务器标志
//	  [eventMsg: string]       — 活动消息
//	  [rateModifier: byte]     — 经验倍率 (写死 100)
//	  [eventXp: byte]          — 活动经验加成 (写死 0)
//	  [rateModifier2: byte]    — 金币倍率 (写死 100)
//	  [dropRate: byte]         — 掉宝率加成 (写死 0)
//	  [reserved: byte(0)]      — 保留
//	  [channelCount: byte]     — 频道数量
//	  // 每个频道:
//	    [chName: string]       — 频道名称 (如 "Scania-1")
//	    [capacity: int]        — 频道容量
//	    [worldID: byte(1)]     — 世界 ID
//	    [channelID: byte]      — 频道序号 (从 0 开始)
//	    [isAdult: bool]        — 是否成人频道
//	  [terminator: short(0)]   — 结束标记
func ServerListEntry(world *WorldInfo) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginServerList)

	w.WriteByte(byte(world.ID))
	w.WriteString(world.Name)
	w.WriteByte(byte(world.Flag))
	w.WriteString(world.EventMessage)
	w.WriteByte(100) // rate modifier (经验倍率)
	w.WriteByte(0)   // event xp
	w.WriteByte(100) // rate modifier (金币倍率)
	w.WriteByte(0)   // drop rate
	w.WriteByte(0)   // reserved
	w.WriteByte(byte(len(world.Channels)))
	for _, ch := range world.Channels {
		w.WriteString(ch.Name)
		w.WriteInt(uint32(ch.Capacity))
		w.WriteByte(ch.WorldID)
		w.WriteByte(ch.ID - 1) // 频道序号从 0 开始
		w.WriteBool(ch.IsAdult)
	}
	w.WriteShort(0) // 结束标记

	return w.Packet()
}

// EndOfServerList 构造服务器列表结束标记 (包7)
//
//	对应 Java: PacketCreator.getEndOfServerList()
//
//	格式: [opcode:0x0A][marker:0xFF]
func EndOfServerList() *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginServerList)
	w.WriteByte(0xFF)
	return w.Packet()
}

// ──────────── LastConnectedWorld (0x1A) ────────────

// LastConnectedWorld 构造"上次连接的世界"封包 (包8)
//
//	对应 Java: PacketCreator.selectWorld()
//
//	格式: [opcode:0x1A][worldID: int(LE)]
func LastConnectedWorld(worldID int32) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginLastConnectedWorld)
	w.WriteInt(uint32(worldID))
	return w.Packet()
}

// ──────────── RecommendedWorlds (0x1B) ────────────

// ──────────── AddNewCharEntry (0x0E) ────────────

// AddNewCharEntry 构造新角色创建成功封包
//
//	对应 Java: PacketCreator.addNewCharEntry()
//
//	格式: [opcode:0x0E][status:byte(0)] + addCharEntry(chr)
func AddNewCharEntry(chr *model.Character) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginAddNewCharEntry)
	w.WriteByte(0) // status = success
	addCharEntry(w, chr)
	return w.Packet()
}

// ──────────── DeleteCharResponse (0x0F) ────────────

// DeleteCharResponse 构造角色操作错误封包（创建/删除失败）
//
//	对应 Java: PacketCreator.deleteCharResponse()
//
//	格式: [opcode:0x0F][charID:int][state:byte]
//	state: 9=未知错误, 10=连接过多, 18=待处理婚礼, 等等
func DeleteCharResponse(charID int32, state byte) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginDeleteCharResponse)
	w.WriteInt(uint32(charID))
	w.WriteByte(state)
	return w.Packet()
}

// ──────────── CharNameResponse (0x0D) ────────────

// CharNameResponse 构造角色名检查结果封包
//
//	对应 Java: PacketCreator.charNameResponse()
//
//	格式: [opcode:0x0D][name:string][nameUsed:bool]
func CharNameResponse(name string, nameUsed bool) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginCharNameResp)
	w.WriteString(name)
	w.WriteBool(nameUsed)
	return w.Packet()
}

// ──────────── ServerStatus (0x03) ────────────

// ServerStatus 构造服务器状态封包（满员等错误）
//
//	对应 Java: PacketCreator.getServerStatus()
//
//	格式: [opcode:0x03][status:short]
func ServerStatus(status int) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginServerStatus)
	w.WriteShort(uint16(status))
	return w.Packet()
}

// ──────────── CharList (0x0B) ────────────

// CharList 构造角色列表封包 (步骤B)
//
//	对应 Java: PacketCreator.getCharList()
//
//	格式:
//	  [opcode:0x0B][status:byte(0)][charCount:byte]
//	  // 每个角色:
//	    addCharStats(...)
//	    addCharLook(...)
//	    [viewall:byte(0)]
//	    // 非GM:
//	    [rankEnabled:byte(1)][rank:int][rankMove:int][jobRank:int][jobRankMove:int]
//	  [picEnabled:byte(2)][charSlots:int]
func CharList(characters []model.Character, charSlots int) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginCharList)

	// status = 0 (成功)
	w.WriteByte(0)
	// 角色数量
	w.WriteByte(byte(len(characters)))

	for i := range characters {
		chr := &characters[i]
		addCharEntry(w, chr)
	}

	// PIC 状态: 2 = 已禁用
	w.WriteByte(2)
	// 可用角色槽位数量
	w.WriteInt(uint32(charSlots))

	return w.Packet()
}

// addCharEntry 写入单个角色的条目信息
//
//	对应 Java: PacketCreator.addCharEntry()
func addCharEntry(w *codec.Writer, chr *model.Character) {
	addCharStats(w, chr)
	addCharLook(w, chr)
	// viewall = false
	w.WriteByte(0)
	// GM 检查：如果是 GM，直接 return (Java: chr.isGM() || chr.isGmJob())
	if chr.Gm > 0 {
		w.WriteByte(0)
		return
	}
	// world rank enabled
	w.WriteByte(1)
	w.WriteInt(uint32(chr.Rank))
	w.WriteInt(uint32(chr.RankMove))
	w.WriteInt(uint32(chr.JobRank))
	w.WriteInt(uint32(chr.JobRankMove))
}

// addCharStats 写入角色属性
//
//	对应 Java: PacketCreator.addCharStats()
func addCharStats(w *codec.Writer, chr *model.Character) {
	w.WriteInt(uint32(chr.ID))        // character id
	w.WritePaddedString(chr.Name, 13) // name (13B padded with \0)
	w.WriteByte(byte(chr.Gender))     // gender
	w.WriteByte(byte(chr.Skincolor))  // skin color
	w.WriteInt(uint32(chr.Face))      // face
	w.WriteInt(uint32(chr.Hair))      // hair
	// pet IDs[3] — 暂无宠物数据，写 0
	for range 3 {
		w.WriteLong(0)
	}
	w.WriteByte(byte(chr.Level))      // level
	w.WriteShort(uint16(chr.Job))     // job
	w.WriteShort(uint16(chr.AttrStr)) // str
	w.WriteShort(uint16(chr.AttrDex)) // dex
	w.WriteShort(uint16(chr.AttrInt)) // int
	w.WriteShort(uint16(chr.AttrLuk)) // luk
	w.WriteShort(uint16(chr.Hp))      // hp
	w.WriteShort(uint16(chr.Maxhp))   // maxhp
	w.WriteShort(uint16(chr.Mp))      // mp
	w.WriteShort(uint16(chr.Maxmp))   // maxmp
	w.WriteShort(uint16(chr.Ap))      // remaining ap
	w.WriteShort(0)                   // remaining sp (简化：后续解析 Sp 字符串)
	w.WriteInt(uint32(chr.Exp))       // exp
	w.WriteShort(uint16(chr.Fame))    // fame
	w.WriteInt(uint32(chr.Gachaexp))  // gacha exp
	w.WriteInt(uint32(chr.Map))       // current map id
	w.WriteByte(byte(chr.Spawnpoint)) // spawnpoint
	w.WriteInt(0)                     // reserved
}

// addCharLook 写入角色外观
//
//	对应 Java: PacketCreator.addCharLook()
func addCharLook(w *codec.Writer, chr *model.Character) {
	w.WriteByte(byte(chr.Gender))    // gender
	w.WriteByte(byte(chr.Skincolor)) // skin color
	w.WriteInt(uint32(chr.Face))     // face
	w.WriteBool(false)               // mega = false
	w.WriteInt(uint32(chr.Hair))     // hair
	// addCharEquips — 暂无装备数据
	// 普通装备结束标记
	w.WriteByte(0xFF)
	// 现金装备结束标记
	w.WriteByte(0xFF)
	// 武器
	w.WriteInt(0)
	// pet item IDs[3]
	for range 3 {
		w.WriteInt(0)
	}
}

// ──────────── ServerIP (0x0C) ────────────

// ServerIP 构造频道服务器IP地址封包 (步骤D)
//
//	对应 Java: PacketCreator.getServerIP()
//
//	格式:
//	  [opcode:0x0C][reserved:short(0)]
//	  [ipAddr:4B][port:short]
//	  [charID:int][reserved:5B(0x00)]
func ServerIP(ip [4]byte, port int, charID int32) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginServerIP)
	w.WriteShort(0)                     // reserved
	w.WriteBytes(ip[:])                 // IP 地址 (4 bytes)
	w.WriteShort(uint16(port))          // 端口
	w.WriteInt(uint32(charID))          // 角色 ID
	w.WriteBytes([]byte{0, 0, 0, 0, 0}) // reserved 5 bytes
	return w.Packet()
}

// ──────────── AfterLoginError (0x09) ────────────

// AfterLoginError 构造角色选择错误响应封包
//
//	对应 Java: PacketCreator.getAfterLoginError()
//
//	格式: [opcode:0x09][reason:short]
func AfterLoginError(reason int) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginSelectCharByVAC)
	w.WriteShort(uint16(reason))
	return w.Packet()
}

// RecommendedWorlds 构造推荐世界消息封包 (包9)
//
//	对应 Java: PacketCreator.sendRecommended()
//
//	格式:
//	  [opcode:0x1B][count:byte]
//	  // 每个推荐世界:
//	    [worldID: int(LE)]
//	    [message: string]
func RecommendedWorlds(list []RecommendedWorld) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginRecommendedWorld)
	w.WriteByte(byte(len(list)))
	for _, rw := range list {
		w.WriteInt(uint32(rw.WorldID))
		w.WriteString(rw.Message)
	}
	return w.Packet()
}
