package model

import (
	"gorm.io/gorm"
)

// Monstercarddata 对应表 monstercarddata
type Monstercarddata struct {
	gorm.Model
	Cardid int `gorm:"column:cardid" json:"cardid"`
	Mobid  int `gorm:"column:mobid" json:"mobid"`
}

func (Monstercarddata) TableName() string {
	return "monstercarddata"
}
