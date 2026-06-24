package channel

import "beidou-go/internal/network"

// ChannelHandler 封包处理器接口
//
// 定义在 login 包而非 handler 包，避免 import cycle:
//
//	login_server.go 使用此接口 → 不依赖 handler 包
//	handler/auth.go 实现此接口 → 可以 import login 包
type ChannelHandler interface {
	HandleCheckPassword(sess *network.Session, username, password string)
	HandleServerList(sess *network.Session)
	HandleServerStatusRequest(sess *network.Session, worldID int)
	HandleCheckCharName(sess *network.Session, name string)
	HandleCharList(sess *network.Session, worldID, channel byte)
	HandleCharCreate(sess *network.Session, data []byte)
	HandleCharSelect(sess *network.Session, charID int32)
}
