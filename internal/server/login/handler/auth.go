package handler

import (
	"time"

	"github.com/sirupsen/logrus"

	"beidou-go/internal/network"
	"beidou-go/internal/network/codec"
	"beidou-go/internal/opcode"
	"beidou-go/internal/server/login"
	"beidou-go/internal/store"
)

// AuthHandler 登录认证处理器
type AuthHandler struct {
	store       *store.AccountStore
	coordinator *login.SessionCoordinator
	autoReg     bool
	log         *logrus.Logger
}

// NewAuthHandler 创建 AuthHandler
func NewAuthHandler(accountStore *store.AccountStore, coord *login.SessionCoordinator, autoReg bool, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		store:       accountStore,
		coordinator: coord,
		autoReg:     autoReg,
		log:         log,
	}
}

// HandleCheckPassword 处理密码验证 (opcode = 0x01 LOGIN_PASSWORD)
//
// 封包格式: [login:string] [password:string] [padding:6B] [hwid:4B]
func (h *AuthHandler) HandleCheckPassword(sess *network.Session, username, password string) {
	// 2. 查询数据库
	account, err := h.store.FindByName(username)
	if err != nil {
		if err == store.ErrAccountNotFound {
			// 账号不存在
			if h.autoReg {
				h.log.Infof("[Auth] 自动注册: account=%s", username)
				hash, hashErr := login.HashPassword(password)
				if hashErr != nil {
					h.log.Errorf("[Auth] bcrypt hash 失败: %v", hashErr)
					sendAndLog(sess, login.LoginStatusFailed(opcode.LoginDBFail), h.log)
					return
				}
				account, err = h.store.Create(username, hash)
				if err != nil {
					h.log.Errorf("[Auth] 创建账号失败: %v", err)
					sendAndLog(sess, login.LoginStatusFailed(opcode.LoginDBFail), h.log)
					return
				}
			} else {
				h.log.Infof("[Auth] 账号不存在: account=%s", username)
				sendAndLog(sess, login.LoginStatusFailed(opcode.LoginNotRegistered), h.log)
				return
			}
		} else {
			h.log.Errorf("[Auth] 数据库查询失败: %v", err)
			sendAndLog(sess, login.LoginStatusFailed(opcode.LoginDBFail), h.log)
			return
		}
	}

	// 3. 封禁检测
	if account.Banned {
		h.log.Infof("[Auth] 账号已封禁: account=%s", username)
		sendAndLog(sess, login.LoginStatusFailed(opcode.LoginBanned), h.log)
		return
	}

	// 4. 密码校验
	if !login.VerifyPassword(password, account.Password) {
		h.log.Infof("[Auth] 密码错误: account=%s", username)
		sendAndLog(sess, login.LoginStatusFailed(opcode.LoginWrongPassword), h.log)
		return
	}

	// 5. 旧 hash 迁移：如果密码是非 bcrypt 格式，自动升级
	if login.NeedsRehash(account.Password) {
		newHash, hashErr := login.HashPassword(password)
		if hashErr != nil {
			h.log.Warnf("[Auth] bcrypt 迁移失败: account=%s, err=%v", username, hashErr)
		} else {
			if err := h.store.UpdatePassword(account.ID, newHash); err != nil {
				h.log.Warnf("[Auth] 密码迁移写入失败: account=%s, err=%v", username, err)
			} else {
				h.log.Infof("[Auth] 密码已迁移到 bcrypt: account=%s", username)
				account.Password = newHash
			}
		}
	}

	// 6. 多开检测
	if err := h.coordinator.AttemptLogin(username, sess.ID()); err != nil {
		h.log.Infof("[Auth] 多开拒绝: account=%s", username)
		sendAndLog(sess, login.LoginStatusFailed(opcode.LoginAlreadyLogin), h.log)
		return
	}

	// 7. 更新最后登录时间
	_ = h.store.UpdateLastLogin(account.ID, time.Now())

	// 8. 登录成功
	sess.AccountID = account.ID

	if err := sess.SendPacket(login.LoginStatusSuccess(account)); err != nil {
		logrus.Errorf("[Auth] 登录失败 err:%v", err)
		return
	}

	h.log.Infof("[Auth] 登录成功: account=%s, id=%d, gm=%d", username, account.ID, 0)
}

// HandleServerList 处理服务器列表请求 (待实现)
func (h *AuthHandler) HandleServerList(sess *network.Session) {
	h.log.Debugf("[Auth] HandleServerList: session=%d (待实现)", sess.ID())
}

// HandleCharList 处理角色列表请求 (待实现)
func (h *AuthHandler) HandleCharList(sess *network.Session, worldID byte) {
	h.log.Debugf("[Auth] HandleCharList: session=%d, world=%d (待实现)", sess.ID(), worldID)
}

// HandleCharCreate 处理角色创建请求 (待实现)
func (h *AuthHandler) HandleCharCreate(sess *network.Session, data []byte) {
	h.log.Debugf("[Auth] HandleCharCreate: session=%d (待实现)", sess.ID())
}

// HandleCharSelect 处理角色选择（进入游戏）(待实现)
func (h *AuthHandler) HandleCharSelect(sess *network.Session, charID int32) {
	h.log.Debugf("[Auth] HandleCharSelect: session=%d, charID=%d (待实现)", sess.ID(), charID)
}

// GetCoordinator 暴露 coordinator，供 login_server 在连接断开时调用 Logout
func (h *AuthHandler) GetCoordinator() *login.SessionCoordinator {
	return h.coordinator
}

// GetAccountStore 暴露 store
func (h *AuthHandler) GetAccountStore() *store.AccountStore {
	return h.store
}

// ──────────────────────────────────────────────
// 内部工具函数
// ──────────────────────────────────────────────

// sendAndLog 发送封包并记录日志
func sendAndLog(sess *network.Session, p *codec.Packet, log *logrus.Logger) {
	if err := sess.SendPacket(p); err != nil {
		log.Errorf("[Auth] 发送封包失败: session=%d, opcode=0x%04X, err=%v", sess.ID(), p.Opcode, err)
	}
}
