package model

import (
	"gorm.io/gorm"
)

// DropDataGlobal 对应表 drop_data_global
type DropDataGlobal struct {
	gorm.Model

	Continent       int    `gorm:"column:continent" json:"continent,omitempty"`
	Itemid          int    `gorm:"column:itemid" json:"itemid,omitempty"`
	MinimumQuantity int    `gorm:"column:minimumQuantity" json:"minimumQuantity,omitempty"`
	MaximumQuantity int    `gorm:"column:maximumQuantity" json:"maximumQuantity,omitempty"`
	Questid         int    `gorm:"column:questid" json:"questid,omitempty"`
	Chance          int    `gorm:"column:chance" json:"chance,omitempty"`
	Comments        string `gorm:"column:comments;type:varchar(200)" json:"comments,omitempty"`
}

func (DropDataGlobal) TableName() string {
	return "drop_data_global"
}
