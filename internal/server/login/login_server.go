package login

import (
	"beidou-go/config"
	"beidou-go/internal/network"

	"github.com/sirupsen/logrus"
)

// Server 登录服务器
type Server struct {
	cfg     *config.Config
	tcpSrv  *network.TCPServer
	log     *logrus.Logger
}

// NewServer 创建登录服务器
func NewServer(cfg *config.Config, tcpSrv *network.TCPServer, log *logrus.Logger) *Server {
	return &Server{
		cfg:    cfg,
		tcpSrv: tcpSrv,
		log:    log,
	}
}

// Start 启动登录服务器，在指定端口上开始监听
func (s *Server) Start() error {
	return s.tcpSrv.Listen(s.cfg.Login.Port, s.handleConnection)
}

// handleConnection 处理客户端连接
func (s *Server) handleConnection(sess *network.Session) {
	defer sess.Close()

	s.log.Infof("[Login] 新连接: %s (session_id=%d)", sess.RemoteAddr(), sess.ID())

	// TODO: 实现握手协议、认证、角色列表等
	// 当前为骨架阶段，连接后直接关闭
	s.log.Debugf("[Login] 连接关闭: session_id=%d", sess.ID())
}
