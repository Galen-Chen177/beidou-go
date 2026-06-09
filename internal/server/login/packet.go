package login

import (
	"beidou-go/internal/model"
	"beidou-go/internal/network/codec"
	"beidou-go/internal/opcode"
)

// ──────────── LoginStatus (0x00) ────────────

// LoginStatusSuccess 构造登录成功响应封包
//
//	格式: [status:0][account_id:4B LE][gender:1B][gm:1B][name:str]
func LoginStatusSuccess(account *model.Account) *codec.Packet {
	w := codec.NewWriterWithOpcode(opcode.LoginStatus)

	w.WriteByte(0)                    // status = success
	w.WriteInt(uint32(account.ID))    // account ID
	w.WriteByte(byte(account.Gender)) // gender
	w.WriteByte(byte(0))              // GM level
	w.WriteString(account.Name)       // account name

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
