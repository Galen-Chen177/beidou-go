package model

import (
	"gorm.io/gorm"
)

// Cooldown 对应表 cooldowns
type Cooldown struct {
	gorm.Model

	Charid    int   `gorm:"column:charid" json:"charid,omitempty"`
	Skillid   int   `gorm:"column:skillid" json:"skillid,omitempty"`
	Length    int64 `gorm:"column:length" json:"length,omitempty"`
	Starttime int64 `gorm:"column:starttime" json:"starttime,omitempty"`
}

func (Cooldown) TableName() string {
	return "cooldowns"
}
