package model

import (
	"gorm.io/gorm"
)

// Savedlocation mapped from table "savedlocations"
type Savedlocation struct {
	gorm.Model
	Characterid  int    `gorm:"column:characterid" json:"characterid,omitempty"`
	Locationtype string `gorm:"column:locationtype;type:varchar(200)" json:"locationtype,omitempty"`
	Map          int    `gorm:"column:map" json:"map,omitempty"`
	Portal       int    `gorm:"column:portal" json:"portal,omitempty"`
}

func (Savedlocation) TableName() string {
	return "savedlocations"
}
