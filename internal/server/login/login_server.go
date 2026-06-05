package login

import (
	"beidou-go/config"
	"beidou-go/internal/network"
	"beidou-go/internal/network/codec"
	"io"

	"github.com/sirupsen/logrus"
)

// Server 登录服务器
type Server struct {
	cfg    *config.Config
	tcpSrv *network.TCPServer
	log    *logrus.Logger
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

// GMS v0.83 服务端握手包（SERVER_HELLO），从 Java 服务端 Wireshark 抓包还原
// 格式：[2字节长度 LE][2字节版本号][2字节子版本][4字节 sendIV][4字节 recvIV][2字节标志]
// 握手阶段使用 2 字节包头（非加密），后续正常通信才切换到 4 字节包头 + AES
var serverHello = []byte{
	0x0E, 0x00, // 包体长度 14 (LE uint16)
	0x53, 0x00, // 版本号 83 (LE uint16)
	0x01, 0x00, // 子版本
	0x31, 0x46, 0x72, 0x7A, // sendIV
	0x60, 0x52, 0x30, 0x78, // recvIV
	0xC1, 0x08, // 服务器标志
}

// handleConnection 处理客户端连接
func (s *Server) handleConnection(sess *network.Session) {
	defer sess.Close()

	s.log.Infof("[Login] 新连接: %s (session_id=%d)", sess.RemoteAddr(), sess.ID())

	// === 第 1 步：发送 SERVER_HELLO ===
	// 冒险岛 v0.83 是服务端先发言，不发这个客户端会一直等
	s.log.Infof("[Login] 发送 SERVER_HELLO: session_id=%d", sess.ID())
	if s.log.IsLevelEnabled(logrus.DebugLevel) {
		s.log.Debugf("[Login] SERVER_HELLO hex:\n%s", codec.HexDump(serverHello))
	}
	if err := sess.Send(serverHello); err != nil {
		s.log.Errorf("[Login] SERVER_HELLO 发送失败: session_id=%d, err=%v", sess.ID(), err)
		return
	}

	// === 第 2 步：接收 CLIENT_HELLO ===
	// 客户端收到 SERVER_HELLO 后会回应 CLIENT_HELLO（2 字节包头格式）
	// 从 CLIENT_HELLO 中可以提取客户端的 IV，用于后续 AES 加解密
	clientHello, err := readHandshakePacket(sess, s.log)
	if err != nil {
		if err == io.EOF {
			s.log.Infof("[Login] 客户端在握手阶段断开: session_id=%d", sess.ID())
		} else {
			s.log.Errorf("[Login] 读取 CLIENT_HELLO 失败: session_id=%d, err=%v", sess.ID(), err)
		}
		return
	}
	s.log.Infof("[Login] 收到 CLIENT_HELLO: session_id=%d, len=%d", sess.ID(), len(clientHello))

	// TODO: 从 CLIENT_HELLO 提取客户端 IV，初始化 MapleCrypto
	// 后续正常通信使用 sess.ReadPacket() / sess.SendPacket()

	// === 第 3 步：循环读取后续封包（正常加密通信阶段） ===
	for {
		packet, err := sess.ReadPacket()
		if err != nil {
			if err == io.EOF {
				s.log.Infof("[Login] 客户端断开: session_id=%d", sess.ID())
			} else {
				s.log.Errorf("[Login] 读取封包失败: session_id=%d, err=%v", sess.ID(), err)
			}
			return
		}

		s.log.Infof("[Login] 收到封包: session_id=%d, opcode=0x%04X, dataLen=%d",
			sess.ID(), packet.Opcode, len(packet.Data))

		// TODO: 根据 opcode 分发到具体 handler
	}
}

// readHandshakePacket 读取握手阶段的 CLIENT_HELLO（无包头，不定长，不加密）
//
// 用 bufio 逐字节试探读取，读到 2 字节后检查是否为合法长度头；
// 如果不是，继续读到超时或封包结束，把所有字节原样返回。
func readHandshakePacket(sess *network.Session, log *logrus.Logger) ([]byte, error) {
	// 先用一个较大的缓冲一次性读（CLIENT_HELLO 通常只有几字节，TCP 会整包送达）
	buf := make([]byte, 512)
	n, err := sess.Read(buf)
	if err != nil {
		return nil, err
	}

	data := buf[:n]
	log.Infof("[Login] CLIENT_HELLO raw (%d bytes):\n%s", n, codec.HexDump(data))

	return data, nil
}
