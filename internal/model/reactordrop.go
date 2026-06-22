package model

import (
	"gorm.io/gorm"
)

// Reactordrop mapped from table "reactordrops"
type Reactordrop struct {
	gorm.Model
	Reactorid int `gorm:"column:reactorid" json:"reactorid,omitempty"`
	Itemid    int `gorm:"column:itemid" json:"itemid,omitempty"`
	Chance    int `gorm:"column:chance" json:"chance,omitempty"`
	Questid   int `gorm:"column:questid" json:"questid,omitempty"`
}

func (Reactordrop) TableName() string {
	return "reactordrops"
}
