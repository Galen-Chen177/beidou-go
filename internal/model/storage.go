package model

import (
	"gorm.io/gorm"
)

// Storage mapped from table "storages"
type Storage struct {
	gorm.Model
	Accountid int `gorm:"column:accountid" json:"accountid,omitempty"`
	World     int `gorm:"column:world" json:"world,omitempty"`
	Slots     int `gorm:"column:slots" json:"slots,omitempty"`
	Meso      int `gorm:"column:meso" json:"meso,omitempty"`
}

func (Storage) TableName() string {
	return "storages"
}
