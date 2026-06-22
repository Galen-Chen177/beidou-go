package model

import (
	"gorm.io/gorm"
)

// Macban 对应表 macbans
type Macban struct {
	gorm.Model
	Macbanid int64  `gorm:"column:macbanid" json:"macbanid"`
	Mac      string `gorm:"column:mac;type:varchar(200)" json:"mac"`
	Aid      string `gorm:"column:aid;type:varchar(200)" json:"aid"`
}

func (Macban) TableName() string {
	return "macbans"
}
