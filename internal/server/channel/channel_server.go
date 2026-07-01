package channel

import (
	"beidou-go/config"
	"beidou-go/internal/crypto"
	"beidou-go/internal/network"
	"beidou-go/internal/network/codec"
	"beidou-go/internal/opcode"
	"encoding/hex"
	"io"

	"github.com/sirupsen/logrus"
)

var version = uint16(83) // 冒险岛 v83

// Server 频道服务器
type Server struct {
	cfg     *config.Config
	tcpSrv  *network.TCPServer
	log     *logrus.Logger
	handler ChannelHandler
}

// NewServer 创建频道服务器
func NewServer(cfg *config.Config, tcpSrv *network.TCPServer, log *logrus.Logger, handler ChannelHandler) *Server {
	return &Server{
		cfg:     cfg,
		tcpSrv:  tcpSrv,
		log:     log,
		handler: handler,
	}
}

// Start 启动频道服务器，在指定端口上开始监听
func (s *Server) Start() error {
	return s.tcpSrv.Listen(s.cfg.Channel.Port, s.handleConnection)
}

// handleConnection 处理客户端连接
func (s *Server) handleConnection(sess *network.Session) {
	defer sess.Close()

	s.log.Infof("[channel] 新连接: %s (session_id=%d)", sess.RemoteAddr(), sess.ID())

	// 1. 构建属于Session的加密解密器
	myCrypto := crypto.NewCrypto(version)
	sess.SetCrypto(myCrypto)

	// 2. 构建并发送SERVER_HELLO
	// 冒险岛 v0.83 是服务端先发言，不发这个客户端会一直等
	serverHello := myCrypto.GenServerHello()
	s.log.Infof("[ServerHello] session_id=%d hex:\n%s", sess.ID(), codec.HexDump(serverHello))
	if err := sess.Send(serverHello); err != nil {
		s.log.Errorf("[channel] SERVER_HELLO 发送失败: session_id=%d, err=%v", sess.ID(), err)
		return
	}

	// 4. 正常通信（首个包为 CLIENT_HELLO 0x0023，解密会使 recvCipher 的 IV 同步到正确状态）
	for {
		packet, err := sess.ReadPacket()
		if err != nil {
			if err == io.EOF {
				s.log.Infof("[channel] 客户端断开: session_id=%d", sess.ID())
			} else {
				s.log.Errorf("[channel] 读取封包失败: session_id=%d, err=%v", sess.ID(), err)
			}
			return
		}
		s.log.Debugf("[Sess %d] === RECV raw (%d bytes) ===\n%s", sess.ID(), len(packet), codec.HexDump(packet))
		decrypted, err := myCrypto.Decrypt(packet)
		if err != nil {
			s.log.Errorf("[channel] 解密失败: session_id=%d, err=%v", sess.ID(), err)
			return
		}
		s.log.Debugf("[Sess %d] === RECV decrypted (%d bytes) ===\n%s", sess.ID(), len(decrypted), hex.Dump(decrypted))
		opcode := uint16(decrypted[0]) | uint16(decrypted[1])<<8
		s.log.Infof("[channel] opcode=0x%04X session_id=%d", opcode, sess.ID())
		data := decrypted[2:]

		// 根据 opcode 分发到具体 handler
		s.dispatch(sess, &codec.Packet{
			Opcode: opcode,
			Data:   data,
		})
	}
}

func (s *Server) dispatch(sess *network.Session, packet *codec.Packet) {
	switch packet.Opcode {
	case opcode.LoginClientHello: // CLIENT_HELLO (0x23)，握手阶段 IV 同步完成
		s.log.Infof("[channel] client hello handshake complete: session_id=%d", sess.ID())

	case opcode.ChannelHello: // PLAYER_LOGGEDIN (0x14)，角色进入频道
		if len(packet.Data) >= 4 {
			charID := int32(packet.Data[0]) | int32(packet.Data[1])<<8 |
				int32(packet.Data[2])<<16 | int32(packet.Data[3])<<24
			s.handler.HandlePlayerLoggedin(sess, charID)
		} else {
			s.log.Warnf("[channel] PLAYER_LOGGEDIN 数据过短: session_id=%d, len=%d", sess.ID(), len(packet.Data))
		}

	default:
		s.log.Warnf("[channel] 未处理的 opcode: session_id=%d, opcode=0x%04X", sess.ID(), packet.Opcode)
	}
}
