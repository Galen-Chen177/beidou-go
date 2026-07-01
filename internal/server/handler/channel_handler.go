package handler

import (
	"github.com/sirupsen/logrus"

	"beidou-go/config"
	"beidou-go/internal/network"
	"beidou-go/internal/server/server_lib"
	"beidou-go/internal/store"
)

// ChannelHandlerImpl 频道服务器封包处理器
//
//	实现 channel.ChannelHandler 接口
type ChannelHandlerImpl struct {
	characterStore *store.CharacterStore
	cfg            *config.Config
	log            *logrus.Logger
}

// NewChannelHandler 创建 ChannelHandlerImpl
func NewChannelHandler(
	characterStore *store.CharacterStore,
	cfg *config.Config,
	log *logrus.Logger,
) *ChannelHandlerImpl {
	return &ChannelHandlerImpl{
		characterStore: characterStore,
		cfg:            cfg,
		log:            log,
	}
}

// HandlePlayerLoggedin 处理 0x14 PLAYER_LOGGEDIN
//
//	对端 Java: PlayerLoggedinHandler.handlePacket()
//
//	客户端发送格式: [charID:int(LE)]
//
//	流程:
//	  1. 从数据库加载角色数据
//	  2. 设置 session 关联的账号/角色/世界信息
//	  3. 发送 GetCharInfo (SET_FIELD 0x7D) — 客户端收到此包后显示游戏画面
func (h *ChannelHandlerImpl) HandlePlayerLoggedin(sess *network.Session, charID int32) {
	h.log.Infof("[Channel] HandlePlayerLoggedin: session=%d, charID=%d", sess.ID(), charID)

	// 1. 从数据库加载角色
	chr, err := h.characterStore.FindByID(charID)
	if err != nil {
		h.log.Warnf("[Channel] 角色不存在: charID=%d, err=%v", charID, err)
		sess.SendPacket(server_lib.AfterLoginError(17))
		return
	}

	// 2. 设置 session 关联数据
	sess.CharID = charID
	sess.AccountID = uint(chr.Accountid)
	sess.WorldID = byte(chr.World)

	// 3. 确定频道号
	channelID := byte(1)
	if sess.ChannelID > 0 {
		channelID = sess.ChannelID
	}

	// 4. 发送 GetCharInfo (SET_FIELD) — 这是最关键的封包
	if err := sess.SendPacket(server_lib.GetCharInfo(chr, channelID)); err != nil {
		h.log.Errorf("[Channel] 发送 GetCharInfo 失败: %v", err)
		return
	}

	h.log.Infof("[Channel] 角色进入游戏成功: name=%s, id=%d, map=%d", chr.Name, chr.ID, chr.Map)
}
