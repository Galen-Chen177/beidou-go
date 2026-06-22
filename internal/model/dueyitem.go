package model

import (
	"gorm.io/gorm"
)

// Dueyitem 对应表 dueyitems
type Dueyitem struct {
	gorm.Model

	Packageid       int64 `gorm:"column:packageid" json:"packageid,omitempty"`
	Inventoryitemid int64 `gorm:"column:inventoryitemid" json:"inventoryitemid,omitempty"`
}

func (Dueyitem) TableName() string {
	return "dueyitems"
}
