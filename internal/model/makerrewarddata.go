package model

import (
	"gorm.io/gorm"
)

// Makerrewarddata 对应表 makerrewarddata (复合主键：itemid + rewardid)
type Makerrewarddata struct {
	gorm.Model
	Itemid   int `gorm:"column:itemid" json:"itemid"`
	Rewardid int `gorm:"column:rewardid" json:"rewardid"`
	Quantity int `gorm:"column:quantity" json:"quantity"`
	Prob     int `gorm:"column:prob" json:"prob"`
}

func (Makerrewarddata) TableName() string {
	return "makerrewarddata"
}
