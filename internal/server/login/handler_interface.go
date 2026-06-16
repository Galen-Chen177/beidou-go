package login

import "beidou-go/internal/network"

// PacketHandler 封包处理器接口
//
// 定义在 login 包而非 handler 包，避免 import cycle:
//
//	login_server.go 使用此接口 → 不依赖 handler 包
//	handler/auth.go 实现此接口 → 可以 import login 包
type PacketHandler interface {
	HandleCheckPassword(sess *network.Session, username, password string)
	HandleServerList(sess *network.Session)
	HandleCharList(sess *network.Session, worldID byte)
	HandleCharCreate(sess *network.Session, data []byte)
	HandleCharSelect(sess *network.Session, charID int32)
}
