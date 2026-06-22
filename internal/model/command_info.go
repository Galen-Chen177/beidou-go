package model

import (
	"gorm.io/gorm"
)

// CommandInfo 对应表 command_info
type CommandInfo struct {
	gorm.Model

	Level        int    `gorm:"column:level" json:"level,omitempty"`
	Syntax       string `gorm:"column:syntax;type:varchar(200)" json:"syntax,omitempty"`
	DefaultLevel int    `gorm:"column:defaultLevel" json:"defaultLevel,omitempty"`
	Clazz        string `gorm:"column:clazz;type:varchar(200)" json:"clazz,omitempty"`
	Enabled      bool   `gorm:"column:enabled" json:"enabled,omitempty"`
}

func (CommandInfo) TableName() string {
	return "command_info"
}
