package login

import (
	"errors"
	"sync"
)

// ErrAlreadyLoggedIn 账号已在其他位置登录
var ErrAlreadyLoggedIn = errors.New("already logged in")

// SessionCoordinator 内存级多开检测
//
// 记录当前在线账号，防止同一账号重复登录。
type SessionCoordinator struct {
	mu      sync.RWMutex
	online  map[string]uint32 // accountName → sessionID
}

// NewSessionCoordinator 创建多开检测器
func NewSessionCoordinator() *SessionCoordinator {
	return &SessionCoordinator{
		online: make(map[string]uint32),
	}
}

// AttemptLogin 尝试登录
// 如果账号已在线，返回 ErrAlreadyLoggedIn
func (sc *SessionCoordinator) AttemptLogin(name string, sessionID uint32) error {
	sc.mu.Lock()
	defer sc.mu.Unlock()

	if _, exists := sc.online[name]; exists {
		return ErrAlreadyLoggedIn
	}
	sc.online[name] = sessionID
	return nil
}

// Logout 登出，移除在线记录
func (sc *SessionCoordinator) Logout(name string) {
	sc.mu.Lock()
	defer sc.mu.Unlock()
	delete(sc.online, name)
}

// OnlineCount 返回当前在线账号数
func (sc *SessionCoordinator) OnlineCount() int {
	sc.mu.RLock()
	defer sc.mu.RUnlock()
	return len(sc.online)
}
