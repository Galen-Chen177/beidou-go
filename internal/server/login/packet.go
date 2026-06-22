package login

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
	w.WriteInt(uint32(account.ID))     // 账号 ID
	w.WriteByte(byte(account.Gender))  // 性别 (0=男, 1=女)

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
	w.WriteByte(0)   // quietBan = false
	w.WriteLong(0)   // quietBan timestamp

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

