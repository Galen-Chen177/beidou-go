package model

import (
	"gorm.io/gorm"
)

// Macfilter 对应表 macfilters
type Macfilter struct {
	gorm.Model
	Macfilterid int64  `gorm:"column:macfilterid" json:"macfilterid"`
	Filter      string `gorm:"column:filter;type:varchar(200)" json:"filter"`
}

func (Macfilter) TableName() string {
	return "macfilters"
}
