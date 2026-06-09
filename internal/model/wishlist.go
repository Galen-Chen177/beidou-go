package model

import (
	"time"

	"gorm.io/gorm"
)

// Wishlist mapped from table "wishlists"
type Wishlist struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ListID    int            `gorm:"column:id;autoIncrement" json:"listID,omitempty"`
	Charid    int            `gorm:"column:charid" json:"charid,omitempty"`
	Sn        int            `gorm:"column:sn" json:"sn,omitempty"`
}

func (Wishlist) TableName() string {
	return "wishlists"
}
