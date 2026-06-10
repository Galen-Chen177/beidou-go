package model

import (
	"time"

	"gorm.io/gorm"
)

// CommandInfo 对应表 command_info
type CommandInfo struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Level        int    `gorm:"column:level" json:"level,omitempty"`
	Syntax       string `gorm:"column:syntax" json:"syntax,omitempty"`
	DefaultLevel int    `gorm:"column:defaultLevel" json:"defaultLevel,omitempty"`
	Clazz        string `gorm:"column:clazz" json:"clazz,omitempty"`
	Enabled      bool   `gorm:"column:enabled" json:"enabled,omitempty"`
}

func (CommandInfo) TableName() string {
	return "command_info"
}
