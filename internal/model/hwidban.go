package model

import "gorm.io/gorm"

// Hwidban 实体映射
type Hwidban struct {
	gorm.Model
	Hwidbanid int64  `gorm:"column:hwidbanid" json:"hwidbanid,omitempty"`
	Hwid      string `gorm:"column:hwid;type:varchar(200)" json:"hwid,omitempty"`
}

func (Hwidban) TableName() string {
	return "hwidbans"
}
