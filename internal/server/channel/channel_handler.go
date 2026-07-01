package channel

import "beidou-go/internal/network"

// ChannelHandler 频道服务器封包处理器接口
//
// 定义在 channel 包而非 handler 包，避免 import cycle:
//
//	channel_server.go 使用此接口 → 不依赖 handler 包
//	handler/channel_handler.go 实现此接口 → 可以 import channel 包
type ChannelHandler interface {
	HandlePlayerLoggedin(sess *network.Session, charID int32)
}
