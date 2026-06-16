package login

import (
	"encoding/hex"
	"fmt"
	"io"

	"github.com/sirupsen/logrus"

	"beidou-go/config"
	"beidou-go/internal/crypto"
	"beidou-go/internal/network"
	"beidou-go/internal/network/codec"
	"beidou-go/internal/opcode"
)

var version = uint16(83) // 冒险岛 v83

// Server 登录服务器
type Server struct {
	cfg     *config.Config
	tcpSrv  *network.TCPServer
	log     *logrus.Logger
	handler PacketHandler
}

// NewServer 创建登录服务器
func NewServer(cfg *config.Config, tcpSrv *network.TCPServer, log *logrus.Logger, handler PacketHandler) *Server {
	return &Server{
		cfg:     cfg,
		tcpSrv:  tcpSrv,
		log:     log,
		handler: handler,
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

	// 1. 构建属于Session的加密解密器
	myCrypto := crypto.NewCrypto(version)
	sess.SetCrypto(myCrypto)

	// 2. 构建并发送SERVER_HELLO
	// 冒险岛 v0.83 是服务端先发言，不发这个客户端会一直等
	serverHello := myCrypto.GenServerHello()
	logrus.Infof("[ServerHello] session_id=%d hex:\n%s", sess.ID(), codec.HexDump(serverHello))
	if err := sess.Send(serverHello); err != nil {
		s.log.Errorf("[Login] SERVER_HELLO 发送失败: session_id=%d, err=%v", sess.ID(), err)
		return
	}

	// 4. 接收并解密 CLIENT_HELLO（加密包，解密会使 recvCipher 的 IV 同步到正确状态）
	clientHello, err := sess.ReadPacket()
	if err != nil {
		if err == io.EOF {
			s.log.Infof("[Login] 客户端在握手阶段断开: session_id=%d", sess.ID())
		} else {
			s.log.Errorf("[Login] 读取 CLIENT_HELLO 失败: session_id=%d, err=%v", sess.ID(), err)
		}
		return
	}
	s.log.Infof("[Login] 收到 CLIENT_HELLO (加密): session_id=%d, len=%d hex:\n%s", sess.ID(), len(clientHello), codec.HexDump(clientHello))

	decryptedHello, err := myCrypto.Decrypt(clientHello)
	if err != nil {
		s.log.Errorf("[Login] CLIENT_HELLO 解密失败: session_id=%d, err=%v", sess.ID(), err)
		return
	}
	s.log.Infof("[Login] CLIENT_HELLO 解密后 (%d bytes):\n%s", len(decryptedHello), hex.Dump(decryptedHello))

	// 5. 正常通信
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
		s.log.Debugf("[Sess %d] === RECV raw (%d bytes) ===\n%s", sess.ID(), len(packet), codec.HexDump(packet))
		decrypted, err := myCrypto.Decrypt(packet)
		if err != nil {
			s.log.Errorf("[Login] 解密失败: session_id=%d, err=%v", sess.ID(), err)
			return
		}
		s.log.Debugf("[Sess %d] === RECV decrypted (%d bytes) ===\n%s", sess.ID(), len(decrypted), hex.Dump(decrypted))
		opcode := uint16(decrypted[0]) | uint16(decrypted[1])<<8
		s.log.Infof("[Login] opcode=0x%04X session_id=%d", opcode, sess.ID())
		data := decrypted[2:]

		// 根据 opcode 分发到具体 handler
		s.dispatch(sess, &codec.Packet{
			Opcode: opcode,
			Data:   data,
		})
	}
}

// dispatch 根据 opcode 路由到对应的 handler
// 在这里先解析成能看的懂的数据，然后调用具体的业务逻辑
func (s *Server) dispatch(sess *network.Session, packet *codec.Packet) {
	switch packet.Opcode {
	// 密码验证 (0x01)
	case opcode.LoginCheckPassword:
		pos := 0
		username := ""
		password := ""

		if len(packet.Data) >= 2 {
			nameLen := int(packet.Data[0]) | int(packet.Data[1])<<8
			fmt.Printf("账号名长度: %d\n", nameLen)
			pos += 2
			if pos+nameLen <= len(packet.Data) {
				fmt.Printf("账号名: %q\n", string(packet.Data[pos:pos+nameLen]))
				username = string(packet.Data[pos : pos+nameLen])
				pos += nameLen
			}
		}

		if pos+2 <= len(packet.Data) {
			pwdLen := int(packet.Data[pos]) | int(packet.Data[pos+1])<<8
			fmt.Printf("密码长度: %d\n", pwdLen)
			pos += 2
			if pos+pwdLen <= len(packet.Data) {
				fmt.Printf("密码: %q\n", string(packet.Data[pos:pos+pwdLen]))
				password = string(packet.Data[pos : pos+pwdLen])
				pos += pwdLen
			}
		}

		s.handler.HandleCheckPassword(sess, username, password)

	case opcode.LoginServerListRereq:
		// 重新请求服务器列表 (0x04)
		s.handler.HandleServerList(sess)

	case opcode.LoginCharListReq:
		// 请求角色列表 (0x05)
		if len(packet.Data) >= 1 {
			s.handler.HandleCharList(sess, packet.Data[0])
		}

	case opcode.LoginCharSelect:
		// 选择角色进入游戏 (0x13)
		if len(packet.Data) >= 4 {
			charID := int32(packet.Data[0]) | int32(packet.Data[1])<<8 |
				int32(packet.Data[2])<<16 | int32(packet.Data[3])<<24
			s.handler.HandleCharSelect(sess, charID)
		}

	case opcode.LoginCharCreate:
		// 创建角色 (0x16)
		s.handler.HandleCharCreate(sess, packet.Data)

	default:
		s.log.Warnf("[Login] 未处理的 opcode: session_id=%d, opcode=0x%04X", sess.ID(), packet.Opcode)
	}
}
