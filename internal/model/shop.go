package model

import (
	"gorm.io/gorm"
)

// Shop mapped from table "shops"
type Shop struct {
	gorm.Model
	Npcid int `gorm:"column:npcid" json:"npcid,omitempty"`
}

func (Shop) TableName() string {
	return "shops"
}
