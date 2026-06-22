package model

import (
	"gorm.io/gorm"
)

// Medalmap 对应表 medalmaps
type Medalmap struct {
	gorm.Model
	Characterid   int   `gorm:"column:characterid" json:"characterid"`
	Queststatusid int64 `gorm:"column:queststatusid" json:"queststatusid"`
	Mapid         int   `gorm:"column:mapid" json:"mapid"`
}

func (Medalmap) TableName() string {
	return "medalmaps"
}
