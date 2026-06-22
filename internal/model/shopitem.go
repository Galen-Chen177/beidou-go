package model

import (
	"gorm.io/gorm"
)

// Shopitem mapped from table "shopitems"
type Shopitem struct {
	gorm.Model
	Shopid   int64 `gorm:"column:shopid" json:"shopid,omitempty"`
	Itemid   int   `gorm:"column:itemid" json:"itemid,omitempty"`
	Price    int   `gorm:"column:price" json:"price,omitempty"`
	Pitch    int   `gorm:"column:pitch" json:"pitch,omitempty"`
	Position int   `gorm:"column:position" json:"position,omitempty"`
}

func (Shopitem) TableName() string {
	return "shopitems"
}
