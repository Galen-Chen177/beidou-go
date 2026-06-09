package model

import (
	"time"

	"gorm.io/gorm"
)

// Keymap 对应表 keymap
type Keymap struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Characterid int            `gorm:"column:characterid" json:"characterid"`
	Key         int            `gorm:"column:key" json:"key"`
	Type        int            `gorm:"column:type" json:"type"`
	Action      int            `gorm:"column:action" json:"action"`
}

func (Keymap) TableName() string {
	return "keymap"
}
