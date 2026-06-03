package network

import (
	"fmt"
	"net"
	"sync"
	"sync/atomic"
)

// TCPServer TCP 服务器，管理多个连接的监听和分发
//
// 不同的逻辑服务（Login、Channel）可以注册不同的端口到同一个 TCPServer。
type TCPServer struct {
	host     string
	mu       sync.RWMutex
	sessions map[uint32]*Session
	nextID   atomic.Uint32

	listeners map[int]net.Listener // port -> listener
	handlers  map[int]SessionHandler

	shutdown bool
	quit     chan struct{}
}

// SessionHandler 客户端连接回调
type SessionHandler func(sess *Session)

// NewTCPServer 创建 TCP 服务器
func NewTCPServer(host string) *TCPServer {
	return &TCPServer{
		host:      host,
		sessions:  make(map[uint32]*Session),
		listeners: make(map[int]net.Listener),
		handlers:  make(map[int]SessionHandler),
		quit:      make(chan struct{}),
	}
}

// Listen 在指定端口上开始监听
// handler: 新连接的回调函数
func (s *TCPServer) Listen(port int, handler SessionHandler) error {
	addr := fmt.Sprintf("%s:%d", s.host, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return fmt.Errorf("监听 %s 失败: %w", addr, err)
	}

	s.mu.Lock()
	s.listeners[port] = ln
	s.handlers[port] = handler
	s.mu.Unlock()

	// 启动 accept 循环
	go s.acceptLoop(ln, port, handler)
	return nil
}

// acceptLoop 循环 accept 新连接
func (s *TCPServer) acceptLoop(ln net.Listener, port int, handler SessionHandler) {
	for {
		conn, err := ln.Accept()
		if err != nil {
			s.mu.RLock()
			down := s.shutdown
			s.mu.RUnlock()
			if down {
				return
			}
			continue
		}

		// 创建 session
		id := s.nextID.Add(1)
		sess := newSession(id, conn, s)

		s.mu.Lock()
		s.sessions[id] = sess
		s.mu.Unlock()

		// 每个连接一个 goroutine
		go handler(sess)
	}
}

// removeSession 移除会话（由 Session.Close 调用）
func (s *TCPServer) removeSession(id uint32) {
	s.mu.Lock()
	delete(s.sessions, id)
	s.mu.Unlock()
}

// SessionCount 返回当前活跃连接数
func (s *TCPServer) SessionCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.sessions)
}

// Shutdown 关闭所有监听和连接
func (s *TCPServer) Shutdown() {
	s.mu.Lock()
	s.shutdown = true
	for port, ln := range s.listeners {
		ln.Close()
		delete(s.listeners, port)
	}
	// 关闭所有会话
	for id, sess := range s.sessions {
		sess.Close()
		delete(s.sessions, id)
	}
	s.mu.Unlock()
	close(s.quit)
}
