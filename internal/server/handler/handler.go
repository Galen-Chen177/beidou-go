package handler

import (
	"time"

	"github.com/sirupsen/logrus"

	"beidou-go/config"
	"beidou-go/internal/model"
	"beidou-go/internal/network"
	"beidou-go/internal/network/codec"
	"beidou-go/internal/opcode"
	"beidou-go/internal/server/login"
	"beidou-go/internal/server/server_lib"
	"beidou-go/internal/store"
)

// AuthHandler 登录认证处理器
type AuthHandler struct {
	store          *store.AccountStore
	characterStore *store.CharacterStore
	coordinator    *login.SessionCoordinator
	worldData      *server_lib.WorldDataProvider
	cfg            *config.Config
	autoReg        bool
	log            *logrus.Logger
}

// NewAuthHandler 创建 AuthHandler
func NewAuthHandler(accountStore *store.AccountStore,
	characterStore *store.CharacterStore,
	coord *login.SessionCoordinator,
	worldData *server_lib.WorldDataProvider,
	cfg *config.Config,
	autoReg bool, log *logrus.Logger) *AuthHandler {
	return &AuthHandler{
		store:          accountStore,
		characterStore: characterStore,
		coordinator:    coord,
		worldData:      worldData,
		cfg:            cfg,
		autoReg:        autoReg,
		log:            log,
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
				hash, hashErr := server_lib.HashPassword(password)
				if hashErr != nil {
					h.log.Errorf("[Auth] bcrypt hash 失败: %v", hashErr)
					sendAndLog(sess, server_lib.LoginStatusFailed(opcode.LoginDBFail), h.log)
					return
				}
				account, err = h.store.Create(username, hash)
				if err != nil {
					h.log.Errorf("[Auth] 创建账号失败: %v", err)
					sendAndLog(sess, server_lib.LoginStatusFailed(opcode.LoginDBFail), h.log)
					return
				}
			} else {
				h.log.Infof("[Auth] 账号不存在: account=%s", username)
				sendAndLog(sess, server_lib.LoginStatusFailed(opcode.LoginNotRegistered), h.log)
				return
			}
		} else {
			h.log.Errorf("[Auth] 数据库查询失败: %v", err)
			sendAndLog(sess, server_lib.LoginStatusFailed(opcode.LoginDBFail), h.log)
			return
		}
	}

	// 3. 封禁检测
	if account.Banned {
		h.log.Infof("[Auth] 账号已封禁: account=%s", username)
		sendAndLog(sess, server_lib.LoginStatusFailed(opcode.LoginBanned), h.log)
		return
	}

	// 4. 密码校验
	if !server_lib.VerifyPassword(password, account.Password) {
		h.log.Infof("[Auth] 密码错误: account=%s", username)
		sendAndLog(sess, server_lib.LoginStatusFailed(opcode.LoginWrongPassword), h.log)
		return
	}

	// 5. 旧 hash 迁移：如果密码是非 bcrypt 格式，自动升级
	if server_lib.NeedsRehash(account.Password) {
		newHash, hashErr := server_lib.HashPassword(password)
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
		sendAndLog(sess, server_lib.LoginStatusFailed(opcode.LoginAlreadyLogin), h.log)
		return
	}

	// 7. 更新最后登录时间
	_ = h.store.UpdateLastLogin(account.ID, time.Now())

	// 8. 登录成功
	sess.AccountID = account.ID

	if err := sess.SendPacket(server_lib.LoginStatusSuccess(account)); err != nil {
		logrus.Errorf("[Auth] 登录失败 err:%v", err)
		return
	}

	h.log.Infof("[Auth] 登录成功: account=%s, id=%d, gm=%d", username, account.ID, 0)
}

// HandleServerList 处理服务器列表请求 (包5/6/7/8/9)
//
//	对应 Java: ServerlistRequestHandler.handlePacket()
//
//	流程:
//	  1. 遍历所有世界，发送 SERVERLIST 封包 (包6)
//	  2. 发送服务器列表结束标记 (包7)
//	  3. 发送上次连接的世界 (包8)
//	  4. 发送推荐世界消息 (包9)
func (h *AuthHandler) HandleServerList(sess *network.Session) {
	h.log.Infof("[Auth] HandleServerList: session=%d, worlds=%d", sess.ID(), len(h.worldData.Worlds()))

	// 1. 发送每个世界的服务器列表 (包6)
	for _, world := range h.worldData.Worlds() {
		pkt := server_lib.ServerListEntry(&world)
		if err := sess.SendPacket(pkt); err != nil {
			h.log.Errorf("[Auth] 发送服务器列表失败: world=%d, err=%v", world.ID, err)
			return
		}
		h.log.Debugf("[Auth] 服务器列表已发送: world=%d, name=%s, channels=%d", world.ID, world.Name, len(world.Channels))
	}

	// 2. 发送服务器列表结束标记 (包7)
	if err := sess.SendPacket(server_lib.EndOfServerList()); err != nil {
		h.log.Errorf("[Auth] 发送列表结束标记失败: err=%v", err)
		return
	}

	// 3. 发送上次连接的世界 (包8)
	if err := sess.SendPacket(server_lib.LastConnectedWorld(h.worldData.LastConnectedWorldID())); err != nil {
		h.log.Errorf("[Auth] 发送上次连接世界失败: err=%v", err)
		return
	}

	// 4. 发送推荐世界消息 (包9)
	if err := sess.SendPacket(server_lib.RecommendedWorlds(h.worldData.RecommendedWorlds())); err != nil {
		h.log.Errorf("[Auth] 发送推荐世界消息失败: err=%v", err)
		return
	}

	h.log.Infof("[Auth] 服务器列表流程完成: session=%d", sess.ID())
}

// HandleCheckCharName 处理角色名检查请求 (opcode 0x15)
//
//	对端 Java: CheckCharNameHandler.handlePacket()
//
//	客户端发送格式: [name:string]
//	服务端响应格式: CharNameResponse 封包
func (h *AuthHandler) HandleCheckCharName(sess *network.Session, name string) {
	h.log.Infof("[Auth] HandleCheckCharName: session=%d, name=%s", sess.ID(), name)

	// 查询名字是否已被占用
	_, err := h.characterStore.FindByName(name)
	nameUsed := err == nil // 找到记录说明名字已被占用

	if err := sess.SendPacket(server_lib.CharNameResponse(name, nameUsed)); err != nil {
		h.log.Errorf("[Auth] 发送 CharNameResponse 失败: %v", err)
	}
}

// HandleServerStatusRequest 处理服务器状态请求 (opcode 0x06)
//
//	对端 Java: ServerStatusRequestHandler.handlePacket()
//
//	客户端发送格式: [world:short(LE)]
//	服务端响应格式: ServerStatus 封包
//
//	状态值: 0=正常, 1=拥挤, 2=满员
func (h *AuthHandler) HandleServerStatusRequest(sess *network.Session, worldID int) {
	h.log.Debugf("[Auth] HandleServerStatusRequest: session=%d, world=%d", sess.ID(), worldID)

	// 查找世界是否存在
	world := h.worldData.FindWorld(worldID)
	if world == nil {
		if err := sess.SendPacket(server_lib.ServerStatus(2)); err != nil {
			h.log.Errorf("[Auth] 发送 ServerStatus 失败: %v", err)
		}
		return
	}

	// 简化：始终返回 0（正常），不做真实的负载检测
	status := 0
	// TODO: 后续实现 getWorldCapacityStatus()

	if err := sess.SendPacket(server_lib.ServerStatus(status)); err != nil {
		h.log.Errorf("[Auth] 发送 ServerStatus 失败: %v", err)
	}
}

// HandleCharList 处理角色列表请求 (步骤A/B)
//
//	对端 Java: CharlistRequestHandler.handlePacket()
//
//	客户端发送格式: [skip:byte(0)][world:byte][channel:byte(0-based)]
//	服务端响应格式: CharList 封包
func (h *AuthHandler) HandleCharList(sess *network.Session, worldID, channel byte) {
	h.log.Infof("[Auth] HandleCharList: session=%d, world=%d, channel=%d", sess.ID(), worldID, channel)

	// 1. 查找世界是否存在
	world := h.worldData.FindWorld(int(worldID))
	if world == nil {
		h.log.Warnf("[Auth] 角色列表请求世界不存在: session=%d, world=%d", sess.ID(), worldID)
		if err := sess.SendPacket(server_lib.ServerStatus(2)); err != nil {
			h.log.Errorf("[Auth] 发送 ServerStatus 失败: %v", err)
		}
		return
	}

	// 2. 查找频道是否存在
	ch := h.worldData.FindChannel(int(worldID), channel)
	if ch == nil {
		h.log.Warnf("[Auth] 角色列表请求频道不存在: session=%d, world=%d, channel=%d", sess.ID(), worldID, channel)
		if err := sess.SendPacket(server_lib.ServerStatus(2)); err != nil {
			h.log.Errorf("[Auth] 发送 ServerStatus 失败: %v", err)
		}
		return
	}

	// 3. 保存世界/频道到 session
	sess.WorldID = worldID
	sess.ChannelID = ch.ID

	// 4. 查询角色列表
	characters, err := h.characterStore.FindByAccountAndWorld(sess.AccountID, int(worldID))
	if err != nil {
		h.log.Errorf("[Auth] 查询角色列表失败: session=%d, account=%d, err=%v", sess.ID(), sess.AccountID, err)
		if err := sess.SendPacket(server_lib.ServerStatus(2)); err != nil {
			h.log.Errorf("[Auth] 发送 ServerStatus 失败: %v", err)
		}
		return
	}

	// 5. 获取角色槽位数量
	account, err := h.store.FindByID(sess.AccountID)
	charSlots := 3 // 默认 3 个槽位
	if err == nil {
		charSlots = account.Characterslots
	}

	h.log.Infof("[Auth] 角色列表: session=%d, count=%d, slots=%d", sess.ID(), len(characters), charSlots)

	// 6. 发送角色列表
	if err := sess.SendPacket(server_lib.CharList(characters, charSlots)); err != nil {
		h.log.Errorf("[Auth] 发送角色列表失败: session=%d, err=%v", sess.ID(), err)
	}
}

// HandleCharSelect 处理角色选择（进入游戏）(步骤C/D)
//
//	对端 Java: CharSelectedHandler.handlePacket()
//
//	客户端发送格式: [charID:int(LE)][macs:string][hostString:string]
//	服务端响应格式: ServerIP 封包 (含频道服务器IP+端口)
func (h *AuthHandler) HandleCharSelect(sess *network.Session, charID int32) {
	h.log.Infof("[Auth] HandleCharSelect: session=%d, charID=%d", sess.ID(), charID)

	// 1. 查询角色
	character, err := h.characterStore.FindByID(charID)
	if err != nil {
		h.log.Warnf("[Auth] 角色不存在: charID=%d, err=%v", charID, err)
		if sendErr := sess.SendPacket(server_lib.AfterLoginError(17)); sendErr != nil {
			h.log.Errorf("[Auth] 发送 AfterLoginError 失败: %v", sendErr)
		}
		return
	}

	// 2. 验证角色属于当前账号
	if character.Accountid != int(sess.AccountID) {
		h.log.Warnf("[Auth] 角色不属于当前账号: charID=%d, charAccount=%d, sessAccount=%d",
			charID, character.Accountid, sess.AccountID)
		if sendErr := sess.SendPacket(server_lib.AfterLoginError(17)); sendErr != nil {
			h.log.Errorf("[Auth] 发送 AfterLoginError 失败: %v", sendErr)
		}
		return
	}

	// 3. 使用角色关联的世界ID
	worldID := sess.WorldID
	if worldID == 0 {
		worldID = byte(character.World)
	}

	world := h.worldData.FindWorld(int(worldID))
	if world == nil {
		h.log.Warnf("[Auth] 世界不存在: world=%d", worldID)
		if sendErr := sess.SendPacket(server_lib.AfterLoginError(10)); sendErr != nil {
			h.log.Errorf("[Auth] 发送 AfterLoginError 失败: %v", sendErr)
		}
		return
	}

	// 4. 构建频道IP
	ip := parseIP(h.cfg.Server.Host)
	port := h.cfg.Channel.Port

	h.log.Infof("[Auth] 角色选择成功: charID=%d, name=%s, world=%d, channel=%d, ip=%v, port=%d",
		charID, character.Name, worldID, sess.ChannelID, ip, port)

	// 5. 发送频道服务器IP封包
	if err := sess.SendPacket(server_lib.ServerIP(ip, port, charID)); err != nil {
		h.log.Errorf("[Auth] 发送 ServerIP 失败: session=%d, err=%v", sess.ID(), err)
	}
}

// HandleCharCreate 处理角色创建请求 (opcode 0x16)
//
//	对端 Java: CreateCharHandler.handlePacket()
//
//	客户端发送格式:
//	  [name:string][job:int][face:int][hair:int][hairColor:int]
//	  [skinColor:int][top:int][bottom:int][shoes:int][weapon:int][gender:byte]
//
//	服务端响应:
//	  成功 → AddNewCharEntry 封包
//	  失败 → DeleteCharResponse 封包
func (h *AuthHandler) HandleCharCreate(sess *network.Session, data []byte) {
	h.log.Infof("[Auth] HandleCharCreate: session=%d, dataLen=%d", sess.ID(), len(data))

	pos := 0

	// 1. 解析 name
	if pos+2 > len(data) {
		h.log.Warn("[Auth] HandleCharCreate: 数据过短，无法解析 name")
		return
	}
	nameLen := int(data[pos]) | int(data[pos+1])<<8
	pos += 2
	if pos+nameLen > len(data) {
		h.log.Warn("[Auth] HandleCharCreate: name 长度超出数据")
		return
	}
	name := string(data[pos : pos+nameLen])
	pos += nameLen

	h.log.Infof("[Auth] HandleCharCreate: name=%s", name)

	// 2. 校验名字长度
	if len(name) < 4 || len(name) > 12 {
		if err := sess.SendPacket(server_lib.DeleteCharResponse(0, 9)); err != nil {
			h.log.Errorf("[Auth] 发送 DeleteCharResponse 失败: %v", err)
		}
		return
	}

	// 3. 检查名字是否已被占用
	_, err := h.characterStore.FindByName(name)
	if err == nil {
		h.log.Infof("[Auth] 角色名已存在: name=%s", name)
		if err := sess.SendPacket(server_lib.DeleteCharResponse(0, 9)); err != nil {
			h.log.Errorf("[Auth] 发送 DeleteCharResponse 失败: %v", err)
		}
		return
	}

	// 4. 检查角色槽位
	account, err := h.store.FindByID(sess.AccountID)
	if err != nil {
		h.log.Errorf("[Auth] 查询账号失败: %v", err)
		return
	}
	charSlots := account.Characterslots
	if charSlots <= 0 {
		charSlots = 3
	}
	existingChars, err := h.characterStore.FindByAccountAndWorld(sess.AccountID, int(sess.WorldID))
	if err != nil {
		h.log.Errorf("[Auth] 查询角色数量失败: %v", err)
		return
	}
	if len(existingChars) >= charSlots {
		h.log.Infof("[Auth] 角色槽位已满: account=%d, chars=%d, slots=%d", sess.AccountID, len(existingChars), charSlots)
		if err := sess.SendPacket(server_lib.DeleteCharResponse(0, 9)); err != nil {
			h.log.Errorf("[Auth] 发送 DeleteCharResponse 失败: %v", err)
		}
		return
	}

	// 5. 解析剩余字段 (每个 int = 4B + 最后一个 byte)
	if pos+37 > len(data) { // 9*4 + 1 = 37 bytes remaining
		h.log.Warn("[Auth] HandleCharCreate: 数据不足以解析属性字段")
		return
	}
	readIntLE := func() int {
		v := int(data[pos]) | int(data[pos+1])<<8 | int(data[pos+2])<<16 | int(data[pos+3])<<24
		pos += 4
		return v
	}
	job := readIntLE()
	face := readIntLE()
	hair := readIntLE()
	hairColor := readIntLE()
	skinColor := readIntLE()
	top := readIntLE()
	bottom := readIntLE()
	shoes := readIntLE()
	weapon := readIntLE()
	gender := int(data[pos])

	h.log.Infof("[Auth] HandleCharCreate: job=%d, face=%d, hair=%d+%d, skin=%d, top=%d, bottom=%d, shoes=%d, weapon=%d, gender=%d",
		job, face, hair, hairColor, skinColor, top, bottom, shoes, weapon, gender)

	// 6. 创建角色
	now := time.Now()
	character := &model.Character{
		Accountid:       int(sess.AccountID),
		World:           int(sess.WorldID),
		Name:            name,
		Level:           1,
		Job:             job,
		Face:            face,
		Hair:            hair + hairColor,
		Gender:          gender,
		Skincolor:       skinColor,
		Map:             10000, // Mushroom Town
		Spawnpoint:      0,
		AttrStr:         12,
		AttrDex:         5,
		AttrInt:         4,
		AttrLuk:         4,
		Maxhp:           50,
		Maxmp:           5,
		Hp:              50,
		Mp:              5,
		Ap:              0,
		Sp:              "0",
		Fame:            0,
		Exp:             0,
		Gachaexp:        0,
		BuddyCapacity:   20,
		Equipslots:      24,
		Useslots:        24,
		Setupslots:      24,
		Etcslots:        24,
		LastLogoutTime:  now,
		LastExpGainTime: now,
	}

	if err := h.characterStore.Create(character); err != nil {
		h.log.Errorf("[Auth] 创建角色失败: name=%s, err=%v", name, err)
		if err := sess.SendPacket(server_lib.DeleteCharResponse(0, 9)); err != nil {
			h.log.Errorf("[Auth] 发送 DeleteCharResponse 失败: %v", err)
		}
		return
	}

	h.log.Infof("[Auth] 角色创建成功: name=%s, id=%d, account=%d", name, character.ID, sess.AccountID)

	// 7. 发送新角色条目
	if err := sess.SendPacket(server_lib.AddNewCharEntry(character)); err != nil {
		h.log.Errorf("[Auth] 发送 AddNewCharEntry 失败: %v", err)
	}
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

// parseIP 将配置中的 host 字符串解析为 [4]byte
// localhost/0.0.0.0 映射为 127.0.0.1
func parseIP(host string) [4]byte {
	var ip [4]byte
	switch host {
	case "", "localhost", "0.0.0.0":
		ip = [4]byte{127, 0, 0, 1}
	default:
		// 简单解析 a.b.c.d 格式
		var a, b, c, d int
		if n, _ := sscanf(host, "%d.%d.%d.%d", &a, &b, &c, &d); n == 4 {
			ip = [4]byte{byte(a), byte(b), byte(c), byte(d)}
		} else {
			ip = [4]byte{127, 0, 0, 1}
		}
	}
	return ip
}

// sscanf 简易字符串解析，类似 fmt.Sscanf
func sscanf(str, format string, args ...*int) (int, error) {
	return parseDotted(str, args)
}

func parseDotted(s string, args []*int) (int, error) {
	n := 0
	val := 0
	sign := 1
	hasDigit := false
	for i := 0; i <= len(s); i++ {
		c := byte(0)
		if i < len(s) {
			c = s[i]
		}
		if c >= '0' && c <= '9' {
			val = val*10 + int(c-'0')
			hasDigit = true
		} else {
			if hasDigit && n < len(args) {
				*args[n] = val * sign
				n++
				val = 0
				sign = 1
				hasDigit = false
			}
		}
	}
	return n, nil
}

// sendAndLog 发送封包并记录日志
func sendAndLog(sess *network.Session, p *codec.Packet, log *logrus.Logger) {
	if err := sess.SendPacket(p); err != nil {
		log.Errorf("[Auth] 发送封包失败: session=%d, opcode=0x%04X, err=%v", sess.ID(), p.Opcode, err)
	}
}
