package model

import (
	"gorm.io/gorm"
)

// Monsterbook 对应表 monsterbook (复合主键：charid + cardid)
type Monsterbook struct {
	gorm.Model
	Charid int `gorm:"column:charid" json:"charid"`
	Cardid int `gorm:"column:cardid" json:"cardid"`
	Level  int `gorm:"column:level" json:"level"`
}

func (Monsterbook) TableName() string {
	return "monsterbook"
}
