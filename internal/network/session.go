package network

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"beidou-go/internal/crypto"
	"beidou-go/internal/network/codec"
)

// Session 客户端会话，封装一个 TCP 连接
type Session struct {
	id     uint32
	conn   net.Conn
	crypto crypto.Crypto
	server *TCPServer
	mu     sync.Mutex
	closed bool

	// 关联的账号/角色信息（登录后填充）
	AccountID uint64
	CharID    int32
	WorldID   byte
	ChannelID byte
}

// newSession 创建新会话
func newSession(id uint32, conn net.Conn, server *TCPServer) *Session {
	conn.SetReadDeadline(time.Now().Add(30 * time.Second))
	return &Session{
		id:     id,
		conn:   conn,
		server: server,
	}
}

// ID 返回会话 ID
func (s *Session) ID() uint32 {
	return s.id
}

// SetCrypto 设置加解密器（握手完成后调用）
func (s *Session) SetCrypto(c crypto.Crypto) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.crypto = c
	// 关闭读超时，进入正常游戏通信
	s.conn.SetReadDeadline(time.Time{})
}

// Crypto 返回加解密器
func (s *Session) Crypto() crypto.Crypto {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.crypto
}

// RemoteAddr 返回客户端地址
func (s *Session) RemoteAddr() string {
	return s.conn.RemoteAddr().String()
}

// Send 直接发送原始字节（不做加密，用于握手阶段）
func (s *Session) Send(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, err := s.conn.Write(data)
	return err
}

// Read 从连接读取原始字节
func (s *Session) Read(buf []byte) (int, error) {
	return s.conn.Read(buf)
}

// Close 关闭会话
func (s *Session) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	s.closed = true
	s.server.removeSession(s.id)
	return s.conn.Close()
}

// IsClosed 判断会话是否已关闭
func (s *Session) IsClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

// SetDeadline 设置读写超时
func (s *Session) SetDeadline(t time.Time) error {
	return s.conn.SetDeadline(t)
}

// ReadPacket 读取一个完整的包
func (s *Session) ReadPacket() ([]byte, error) {
	// 1. 读 4 字节包头（不加密）
	header := make([]byte, 4)
	if _, err := io.ReadFull(s.conn, header); err != nil {
		return nil, err
	}
	logrus.Infof("ReadPacket header:[%v]", header)
	bodyLen := decodePacketLength(header)
	logrus.Infof("ReadPacket bodylen:[%v]", bodyLen)

	// 2. 读包体（可能处于加密状态）
	body := make([]byte, bodyLen)
	if _, err := io.ReadFull(s.conn, body); err != nil {
		return nil, err
	}

	// 3. 打印原始字节 → 和 Wireshark 抓到的 TCP payload 完全一致
	Log.Debugf("[Sess %d] === RECV raw (%d bytes) ===\n%s", s.id, bodyLen, codec.HexDump(body))
	return append(header, body...), nil
}

func decodePacketLength(header []byte) int {
	hi := int(header[1] ^ header[3]) // Java: header[1] ^ header[3]
	lo := int(header[0] ^ header[2]) // Java: header[0] ^ header[2]
	return (hi << 8) | lo
}

// SendPacket 编码、加密（如果已完成握手）、发送封包，并在 Debug 级别输出 hex dump。
//
//	opcode+data → [SEND plain] → 加密 → [SEND encrypted] → TCP 字节
func (s *Session) SendPacket(p *codec.Packet) error {
	var err error
	// 1. 编码为包体：[opcode:2B][data]
	body := p.Bytes()

	// Info 级别始终输出封包摘要
	Log.Infof("[Sess %d] SEND opcode=0x%04X  dataLen=%d", s.id, p.Opcode, len(body)-2)
	if Log.IsLevelEnabled(logrus.DebugLevel) {
		Log.Debugf("[Sess %d] === SEND plain (opcode=0x%04X, %d bytes) ===\n%s",
			s.id, p.Opcode, len(body), codec.HexDump(body))
	}

	s.mu.Lock()
	defer s.mu.Unlock()

	// 2. 加密
	if s.crypto != nil {
		body, err = s.crypto.Encrypt(body)
		if err != nil {
			return err
		}
		if Log.IsLevelEnabled(logrus.DebugLevel) {
			Log.Debugf("[Sess %d] --- SEND encrypted ---\n%s", s.id, codec.HexDump(body))
		}
	}

	// 3. 组包：[4 字节包头（小端 = bodyLen）][包体]
	out := make([]byte, 4+len(body))
	binary.LittleEndian.PutUint32(out[:4], uint32(len(body)))
	copy(out[4:], body)

	_, err = s.conn.Write(out)
	return err
}
