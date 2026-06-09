package model

import (
	"time"

	"gorm.io/gorm"
)

// Cooldown 对应表 cooldowns
type Cooldown struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Charid    int   `gorm:"column:charid" json:"charid,omitempty"`
	Skillid   int   `gorm:"column:skillid" json:"skillid,omitempty"`
	Length    int64 `gorm:"column:length" json:"length,omitempty"`
	Starttime int64 `gorm:"column:starttime" json:"starttime,omitempty"`
}

func (Cooldown) TableName() string {
	return "cooldowns"
}
