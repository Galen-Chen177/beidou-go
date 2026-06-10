package model

import (
	"time"

	"gorm.io/gorm"
)

// MtsCart 对应表 mts_cart
type MtsCart struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Cid       int            `gorm:"column:cid" json:"cid"`
	Itemid    int            `gorm:"column:itemid" json:"itemid"`
}

func (MtsCart) TableName() string {
	return "mts_cart"
}
