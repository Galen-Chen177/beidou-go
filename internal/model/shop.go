package model

import (
	"time"

	"gorm.io/gorm"
)

// Shop mapped from table "shops"
type Shop struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Shopid    int64          `gorm:"column:shopid;autoIncrement" json:"shopid,omitempty"`
	Npcid     int            `gorm:"column:npcid" json:"npcid,omitempty"`
}

func (Shop) TableName() string {
	return "shops"
}
