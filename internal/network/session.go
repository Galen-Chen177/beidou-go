package network

import (
	"net"
	"sync"
	"time"

	"beidou-go/internal/crypto"
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
