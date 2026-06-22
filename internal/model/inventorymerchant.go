package model

import (
	"gorm.io/gorm"
)

// Inventorymerchant 对应表 inventorymerchant
type Inventorymerchant struct {
	gorm.Model
	Inventorymerchantid int64 `gorm:"column:inventorymerchantid" json:"inventorymerchantid"`
	Inventoryitemid     int64 `gorm:"column:inventoryitemid" json:"inventoryitemid"`
	Characterid         int   `gorm:"column:characterid" json:"characterid"`
	Bundles             int   `gorm:"column:bundles" json:"bundles"`
}

func (Inventorymerchant) TableName() string {
	return "inventorymerchant"
}
