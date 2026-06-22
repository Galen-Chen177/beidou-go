package model

import (
	"gorm.io/gorm"
)

// Makerrecipedata 对应表 makerrecipedata (复合主键：itemid + reqItem)
type Makerrecipedata struct {
	gorm.Model
	Itemid  int `gorm:"column:itemid" json:"itemid"`
	ReqItem int `gorm:"column:reqItem" json:"reqItem"`
	Count   int `gorm:"column:count" json:"count"`
}

func (Makerrecipedata) TableName() string {
	return "makerrecipedata"
}
