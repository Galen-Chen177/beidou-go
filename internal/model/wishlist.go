package model

import (
	"gorm.io/gorm"
)

// Wishlist mapped from table "wishlists"
type Wishlist struct {
	gorm.Model

	Charid int `gorm:"column:charid" json:"charid,omitempty"`
	Sn     int `gorm:"column:sn" json:"sn,omitempty"`
}

func (Wishlist) TableName() string {
	return "wishlists"
}
