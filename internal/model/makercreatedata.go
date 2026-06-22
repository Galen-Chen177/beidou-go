package model

import (
	"gorm.io/gorm"
)

// Makercreatedata 对应表 makercreatedata (复合主键：id + itemid)
type Makercreatedata struct {
	gorm.Model
	MakerID       int `gorm:"column:id" json:"makerId"`
	Itemid        int `gorm:"column:itemid" json:"itemid"`
	ReqLevel      int `gorm:"column:reqLevel" json:"reqLevel"`
	ReqMakerLevel int `gorm:"column:reqMakerLevel" json:"reqMakerLevel"`
	ReqMeso       int `gorm:"column:reqMeso" json:"reqMeso"`
	ReqItem       int `gorm:"column:reqItem" json:"reqItem"`
	ReqEquip      int `gorm:"column:reqEquip" json:"reqEquip"`
	Catalyst      int `gorm:"column:catalyst" json:"catalyst"`
	Quantity      int `gorm:"column:quantity" json:"quantity"`
	Tuc           int `gorm:"column:tuc" json:"tuc"`
}

func (Makercreatedata) TableName() string {
	return "makercreatedata"
}
