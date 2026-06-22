package model

import (
	"gorm.io/gorm"
)

// Specialcashitem mapped from table "specialcashitems"
type Specialcashitem struct {
	gorm.Model
	ItemID   int `gorm:"column:id" json:"itemID,omitempty"`
	Sn       int `gorm:"column:sn" json:"sn,omitempty"`
	Modifier int `gorm:"column:modifier" json:"modifier,omitempty"`
	Info     int `gorm:"column:info" json:"info,omitempty"`
}

func (Specialcashitem) TableName() string {
	return "specialcashitems"
}
