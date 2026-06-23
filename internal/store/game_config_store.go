package store

import (
	"errors"
	"strconv"

	"beidou-go/internal/model"

	"gorm.io/gorm"
)

// ErrConfigNotFound 配置项不存在
var ErrConfigNotFound = errors.New("config not found")

// GameConfigStore game_config 表数据访问层
type GameConfigStore struct {
	db *gorm.DB
}

// NewGameConfigStore 创建 GameConfigStore
func NewGameConfigStore(db *gorm.DB) *GameConfigStore {
	return &GameConfigStore{db: db}
}

// Get 按 configType + configSubType + configCode 查询单条配置的 configValue
func (s *GameConfigStore) Get(configType, configSubType, configCode string) (string, error) {
	var cfg model.GameConfig
	result := s.db.Where("configType = ? AND configSubType = ? AND configCode = ?",
		configType, configSubType, configCode).First(&cfg)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return "", ErrConfigNotFound
		}
		return "", result.Error
	}
	return cfg.ConfigValue, nil
}

// GetInt 查询 int 类型的配置值，不存在时返回默认值
func (s *GameConfigStore) GetInt(configType, configSubType, configCode string, defaultVal int) int {
	val, err := s.Get(configType, configSubType, configCode)
	if err != nil {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return n
}

// GetFloat 查询 float 类型的配置值，不存在时返回默认值
func (s *GameConfigStore) GetFloat(configType, configSubType, configCode string, defaultVal float64) float64 {
	val, err := s.Get(configType, configSubType, configCode)
	if err != nil {
		return defaultVal
	}
	f, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return defaultVal
	}
	return f
}

// GetString 查询 string 类型的配置值，不存在时返回默认值
func (s *GameConfigStore) GetString(configType, configSubType, configCode string, defaultVal string) string {
	val, err := s.Get(configType, configSubType, configCode)
	if err != nil {
		return defaultVal
	}
	return val
}

// ──────────── 世界配置快捷方法 ────────────

// GetWorldInt 读取世界级 int 配置: world.{id}.{key}
func (s *GameConfigStore) GetWorldInt(worldID int, key string, defaultVal int) int {
	return s.GetInt("world", strconv.Itoa(worldID), key, defaultVal)
}

// GetWorldFloat 读取世界级 float 配置
func (s *GameConfigStore) GetWorldFloat(worldID int, key string, defaultVal float64) float64 {
	return s.GetFloat("world", strconv.Itoa(worldID), key, defaultVal)
}

// GetWorldString 读取世界级 string 配置
func (s *GameConfigStore) GetWorldString(worldID int, key string, defaultVal string) string {
	return s.GetString("world", strconv.Itoa(worldID), key, defaultVal)
}

// ──────────── 服务器全局配置快捷方法 ────────────

// GetServerInt 读取服务器全局 int 配置: server.global.{key}
func (s *GameConfigStore) GetServerInt(key string, defaultVal int) int {
	return s.GetInt("server", "global", key, defaultVal)
}
