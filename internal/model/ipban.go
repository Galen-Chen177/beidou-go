package model

import (
	"gorm.io/gorm"
)

// Ipban 对应表 ipbans
type Ipban struct {
	gorm.Model
	Ipbanid int64  `gorm:"column:ipbanid" json:"ipbanid"`
	Ip      string `gorm:"column:ip;type:varchar(200)" json:"ip"`
	Aid     string `gorm:"column:aid;type:varchar(200)" json:"aid"`
}

func (Ipban) TableName() string {
	return "ipbans"
}
