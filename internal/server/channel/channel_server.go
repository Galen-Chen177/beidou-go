package channel

import (
	"beidou-go/config"
	"beidou-go/internal/network"

	"github.com/sirupsen/logrus"
)

// Server 频道服务器
type Server struct {
	cfg    *config.Config
	tcpSrv *network.TCPServer
	log    *logrus.Logger
}

// NewServer 创建频道服务器
func NewServer(cfg *config.Config, tcpSrv *network.TCPServer, log *logrus.Logger) *Server {
	return &Server{
		cfg:    cfg,
		tcpSrv: tcpSrv,
		log:    log,
	}
}

// Start 启动频道服务器，在指定端口上开始监听
func (s *Server) Start() error {
	return s.tcpSrv.Listen(s.cfg.Channel.Port, s.handleConnection)
}

// handleConnection 处理客户端连接
func (s *Server) handleConnection(sess *network.Session) {
	defer sess.Close()

	s.log.Infof("[Channel] 新连接: %s (session_id=%d)", sess.RemoteAddr(), sess.ID())

	// TODO: 实现地图、战斗、聊天等游戏逻辑
	// 当前为骨架阶段
	s.log.Debugf("[Channel] 连接关闭: session_id=%d", sess.ID())
}
