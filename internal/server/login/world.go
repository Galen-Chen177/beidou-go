package login

import (
	"fmt"

	"github.com/sirupsen/logrus"

	"beidou-go/internal/store"
)

// ──────────────────────────────────────────────
// 世界名称常量 (对齐 Java: GameConstants.WORLD_NAMES)
// ──────────────────────────────────────────────

var defaultWorldNames = []string{
	"Scania", "Bera", "Broa", "Windia", "Khaini", "Bellocan",
	"Mardia", "Kradia", "Yellonde", "Demethos", "Galicia",
	"El Nido", "Zenith", "Arcenia", "Kastia", "Judis",
	"Plana", "Kalluna", "Stius", "Croa", "Medere",
}

// ──────────────────────────────────────────────
// 世界/频道 数据结构
// ──────────────────────────────────────────────

// ChannelInfo 频道信息
type ChannelInfo struct {
	Name     string // 频道名称 (如 "Scania-1")
	Capacity int32  // 频道负载容量 (对客端显示用)
	WorldID  byte   // 世界 ID
	ID       byte   // 频道序号 (1-based)
	IsAdult  bool   // 是否成人频道
}

// WorldInfo 世界信息
type WorldInfo struct {
	ID           int           // 世界 ID
	Name         string        // 世界名称 (如 "Scania")
	Flag         int           // 服务器标志
	EventMessage string        // 活动消息
	Channels     []ChannelInfo // 频道列表
}

// RecommendedWorld 推荐世界
type RecommendedWorld struct {
	WorldID int    // 世界 ID
	Message string // 推荐消息文本
}

// ──────────────────────────────────────────────
// WorldDataProvider 世界数据提供者
// ──────────────────────────────────────────────

// WorldDataProvider 从 game_config 表加载世界/频道数据
//
//	对齐 Java: Server.initWorld() + Server.worldRecommendedList()
type WorldDataProvider struct {
	worlds      []WorldInfo
	recommended []RecommendedWorld
}

// NewWorldDataProvider 从 game_config 表加载世界数据
//
// 加载流程 (对齐 Java Server.initWorld()):
//  1. 读取 server.global.WORLDS → 世界数量
//  2. 对每个世界 i:
//     a. 读取 world.{i}.channel_size → 频道数量
//     b. 读取 world.{i}.flag / event_message / recommend_message
//     c. 为每个频道构造 ChannelInfo
//  3. 构造推荐世界列表
//
// 如果 game_config 中没有配置，默认返回 1 个世界 + 2 个频道。
func NewWorldDataProvider(store *store.GameConfigStore, log *logrus.Logger) *WorldDataProvider {
	p := &WorldDataProvider{
		worlds:      make([]WorldInfo, 0),
		recommended: make([]RecommendedWorld, 0),
	}

	// 1. 读取世界数量
	worldCount := store.GetServerInt("WORLDS", 1)
	log.Infof("[WorldData] 世界数量: %d (server.global.WORLDS)", worldCount)

	// 2. 遍历每个世界
	for i := range worldCount {
		channelSize := store.GetWorldInt(i, "channel_size", 2)
		flag := store.GetWorldInt(i, "flag", 0)
		eventMessage := store.GetWorldString(i, "event_message", "")
		recommendMessage := store.GetWorldString(i, "recommend_message", "")
		// 世界名称：优先从配置读取，否则用默认名称表
		worldName := store.GetWorldString(i, "world_name", "")
		if worldName == "" && i < len(defaultWorldNames) {
			worldName = defaultWorldNames[i]
		}
		if worldName == "" {
			worldName = fmt.Sprintf("World-%d", i)
		}

		world := WorldInfo{
			ID:           i,
			Name:         worldName,
			Flag:         flag,
			EventMessage: eventMessage,
			Channels:     make([]ChannelInfo, 0, channelSize),
		}

		// 构造频道列表
		for j := 1; j <= channelSize; j++ {
			ch := ChannelInfo{
				Name:     fmt.Sprintf("%s-%d", worldName, j),
				Capacity: 1200, // 默认频道容量
				WorldID:  byte(i),
				ID:       byte(j),
				IsAdult:  false,
			}
			world.Channels = append(world.Channels, ch)
		}

		p.worlds = append(p.worlds, world)

		// 推荐消息
		if recommendMessage != "" {
			p.recommended = append(p.recommended, RecommendedWorld{
				WorldID: i,
				Message: recommendMessage,
			})
		}

		log.Infof("[WorldData] 世界 %d: name=%s, flag=%d, channels=%d, event=%q, recommend=%q",
			i, worldName, flag, channelSize, eventMessage, recommendMessage)
	}

	// 兜底：如果 game_config 表完全是空的，提供默认世界
	if len(p.worlds) == 0 {
		log.Warn("[WorldData] game_config 无世界配置，使用默认值")
		p.worlds = []WorldInfo{
			{
				ID:           0,
				Name:         "Scania",
				Flag:         0,
				EventMessage: "",
				Channels: []ChannelInfo{
					{Name: "Scania-1", Capacity: 1200, WorldID: 0, ID: 1, IsAdult: false},
					{Name: "Scania-2", Capacity: 1200, WorldID: 0, ID: 2, IsAdult: false},
				},
			},
		}
	}

	return p
}

// Worlds 返回世界列表
func (p *WorldDataProvider) Worlds() []WorldInfo {
	return p.worlds
}

// RecommendedWorlds 返回推荐世界列表
func (p *WorldDataProvider) RecommendedWorlds() []RecommendedWorld {
	return p.recommended
}

// FindWorld 按 ID 查找世界，未找到返回 nil
func (p *WorldDataProvider) FindWorld(worldID int) *WorldInfo {
	for i := range p.worlds {
		if p.worlds[i].ID == worldID {
			return &p.worlds[i]
		}
	}
	return nil
}

// FindChannel 在世界中按频道号查找频道，未找到返回 nil
func (p *WorldDataProvider) FindChannel(worldID int, channelID byte) *ChannelInfo {
	w := p.FindWorld(worldID)
	if w == nil {
		return nil
	}
	for i := range w.Channels {
		if w.Channels[i].ID == channelID {
			return &w.Channels[i]
		}
	}
	return nil
}

// LastConnectedWorldID 返回上次连接的世界 ID
//
// 对齐 Java: PacketCreator.selectWorld(0)，当前固定返回 0
// 后续可实现为从账号数据读取最近活跃世界
func (p *WorldDataProvider) LastConnectedWorldID() int32 {
	return 0
}
