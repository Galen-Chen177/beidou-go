package model

import (
	"gorm.io/gorm"
)

// Keymap 对应表 keymap
type Keymap struct {
	gorm.Model
	Characterid int `gorm:"column:characterid" json:"characterid"`
	Key         int `gorm:"column:key" json:"key"`
	Type        int `gorm:"column:type" json:"type"`
	Action      int `gorm:"column:action" json:"action"`
}

func (Keymap) TableName() string {
	return "keymap"
}
