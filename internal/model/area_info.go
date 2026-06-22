package model

import (
	"gorm.io/gorm"
)

// AreaInfo 对应表 area_info
type AreaInfo struct {
	gorm.Model

	Charid int    `gorm:"column:charid" json:"charid,omitempty"`
	Area   int    `gorm:"column:area" json:"area,omitempty"`
	Info   string `gorm:"column:info;type:varchar(200)" json:"info,omitempty"`
}

func (AreaInfo) TableName() string {
	return "area_info"
}
