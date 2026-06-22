package model

import (
	"gorm.io/gorm"
)

// Makerreagentdata 对应表 makerreagentdata
type Makerreagentdata struct {
	gorm.Model
	Itemid int    `gorm:"column:itemid" json:"itemid"`
	Stat   string `gorm:"column:stat;type:varchar(200)" json:"stat"`
	Value  int    `gorm:"column:value" json:"value"`
}

func (Makerreagentdata) TableName() string {
	return "makerreagentdata"
}
