package handler

import (
	"beidou-go/internal/network"
)

// AuthHandler 登录认证处理器 (待实现)
type AuthHandler struct{}

// HandleCheckPassword 处理密码验证 (待实现)
func (h *AuthHandler) HandleCheckPassword(sess *network.Session, data []byte) {
	// TODO: 实现密码验证逻辑
	// 1. 解析客户端发来的用户名和密码
	// 2. 查询数据库验证
	// 3. 返回 LoginStatus 封包
}

// HandleServerList 处理服务器列表请求 (待实现)
func (h *AuthHandler) HandleServerList(sess *network.Session) {
	// TODO: 返回世界和频道列表
}

// HandleCharList 处理角色列表请求 (待实现)
func (h *AuthHandler) HandleCharList(sess *network.Session, worldID byte) {
	// TODO: 返回该账号在该世界下的所有角色
}

// HandleCharCreate 处理角色创建请求 (待实现)
func (h *AuthHandler) HandleCharCreate(sess *network.Session, data []byte) {
	// TODO: 解析角色名、属性等，写入数据库
}

// HandleCharSelect 处理角色选择（进入游戏）(待实现)
func (h *AuthHandler) HandleCharSelect(sess *network.Session, charID int32) {
	// TODO: 验证角色归属，通知频道服务器
}
