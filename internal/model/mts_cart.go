package model

import (
	"gorm.io/gorm"
)

// MtsCart 对应表 mts_cart
type MtsCart struct {
	gorm.Model
	Cid    int `gorm:"column:cid" json:"cid"`
	Itemid int `gorm:"column:itemid" json:"itemid"`
}

func (MtsCart) TableName() string {
	return "mts_cart"
}
