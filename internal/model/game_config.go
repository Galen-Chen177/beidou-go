package model

import (
	"time"

	"gorm.io/gorm"
)

// GameConfig 游戏参数表 实体映射
type GameConfig struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ConfigType    string         `json:"configType,omitempty"`
	ConfigSubType string         `json:"configSubType,omitempty"`
	ConfigClazz   string         `json:"configClazz,omitempty"`
	ConfigCode    string         `json:"configCode,omitempty"`
	ConfigValue   string         `json:"configValue,omitempty"`
	ConfigDesc    string         `json:"configDesc,omitempty"`
	UpdateTime    time.Time      `json:"updateTime,omitempty"`
}

func (GameConfig) TableName() string {
	return "game_config"
}
