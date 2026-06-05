package network

import (
	"encoding/binary"
	"io"
	"net"
	"sync"
	"time"

	"beidou-go/internal/crypto"
	"beidou-go/internal/network/codec"

	"github.com/sirupsen/logrus"
)

// Session 客户端会话，封装一个 TCP 连接
type Session struct {
	id     uint32
	conn   net.Conn
	crypto *crypto.MapleCrypto
	server *TCPServer
	mu     sync.Mutex
	closed bool

	// 关联的账号/角色信息（登录后填充）
	AccountID int32
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
func (s *Session) SetCrypto(c *crypto.MapleCrypto) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.crypto = c
	// 关闭读超时，进入正常游戏通信
	s.conn.SetReadDeadline(time.Time{})
}

// Crypto 返回加解密器
func (s *Session) Crypto() *crypto.MapleCrypto {
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

// SendEncrypted 加密后发送
func (s *Session) SendEncrypted(data []byte) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.crypto != nil {
		s.crypto.Encrypt(data)
	}
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

// ──────────────────────────────────────────────
// 封包级收发（带 hex dump 日志，对标 Wireshark）
// ──────────────────────────────────────────────

// ReadPacket 读取一个完整封包，自动解密（如果已完成握手），并在 Debug 级别输出 hex dump。
//
//	TCP 字节 → [RECV raw] → 解密 → [RECV decrypted] → opcode + data
//
// 用法：在 handler goroutine 中循环调用，直到 io.EOF。
func (s *Session) ReadPacket() (*codec.Packet, error) {
	// 1. 读 4 字节包头（不加密）
	header := make([]byte, 4)
	if _, err := io.ReadFull(s.conn, header); err != nil {
		return nil, err
	}
	bodyLen := binary.LittleEndian.Uint32(header)

	// 2. 读包体（可能处于加密状态）
	body := make([]byte, bodyLen)
	if _, err := io.ReadFull(s.conn, body); err != nil {
		return nil, err
	}

	// 3. 打印原始字节 → 和 Wireshark 抓到的 TCP payload 完全一致
	if Log.IsLevelEnabled(logrus.DebugLevel) {
		Log.Debugf("[Sess %d] === RECV raw (%d bytes) ===\n%s", s.id, bodyLen, codec.HexDump(body))
	}

	// 4. 解密（AES/OFB 流解密，就地修改）
	if s.crypto != nil {
		s.crypto.Decrypt(body)
		if Log.IsLevelEnabled(logrus.DebugLevel) {
			Log.Debugf("[Sess %d] --- RECV decrypted ---\n%s", s.id, codec.HexDump(body))
		}
	}

	// 5. 解析 opcode（2 字节小端序）
	if len(body) < 2 {
		return &codec.Packet{Opcode: 0, Data: body}, nil
	}
	opcode := binary.LittleEndian.Uint16(body[:2])

	// Info 级别始终输出封包摘要（不需要开 Debug 也能看到）
	Log.Infof("[Sess %d] RECV opcode=0x%04X  dataLen=%d", s.id, opcode, len(body)-2)
	if Log.IsLevelEnabled(logrus.DebugLevel) {
		Log.Debugf("[Sess %d] RECV hex dump:\n%s", s.id, codec.HexDump(body))
	}

	return &codec.Packet{Opcode: opcode, Data: body[2:]}, nil
}

// SendPacket 编码、加密（如果已完成握手）、发送封包，并在 Debug 级别输出 hex dump。
//
//	opcode+data → [SEND plain] → 加密 → [SEND encrypted] → TCP 字节
func (s *Session) SendPacket(p *codec.Packet) error {
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
		s.crypto.Encrypt(body)
		if Log.IsLevelEnabled(logrus.DebugLevel) {
			Log.Debugf("[Sess %d] --- SEND encrypted ---\n%s", s.id, codec.HexDump(body))
		}
	}

	// 3. 组包：[4 字节包头（小端 = bodyLen）][包体]
	out := make([]byte, 4+len(body))
	binary.LittleEndian.PutUint32(out[:4], uint32(len(body)))
	copy(out[4:], body)

	_, err := s.conn.Write(out)
	return err
}
