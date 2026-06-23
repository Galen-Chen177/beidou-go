package model

import (
	"gorm.io/gorm"
)

// GameConfig 游戏参数表 实体映射
type GameConfig struct {
	gorm.Model
	ConfigType    string `gorm:"column:configType;type:varchar(200)" json:"configType,omitempty"`
	ConfigSubType string `gorm:"column:configSubType;type:varchar(200)" json:"configSubType,omitempty"`
	ConfigClazz   string `gorm:"column:configClazz;type:varchar(200)" json:"configClazz,omitempty"`
	ConfigCode    string `gorm:"column:configCode;type:varchar(200)" json:"configCode,omitempty"`
	ConfigValue   string `gorm:"column:configValue;type:varchar(200)" json:"configValue,omitempty"`
	ConfigDesc    string `gorm:"column:configDesc;type:varchar(200)" json:"configDesc,omitempty"`
}

func (GameConfig) TableName() string {
	return "game_config"
}
