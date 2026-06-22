package model

import (
	"gorm.io/gorm"
)

// Trocklocation mapped from table "trocklocations"
type Trocklocation struct {
	gorm.Model
	Characterid int `gorm:"column:characterid" json:"characterid,omitempty"`
	Mapid       int `gorm:"column:mapid" json:"mapid,omitempty"`
	Vip         int `gorm:"column:vip" json:"vip,omitempty"`
}

func (Trocklocation) TableName() string {
	return "trocklocations"
}
