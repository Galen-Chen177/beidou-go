package model

import (
	"time"

	"gorm.io/gorm"
)

// AreaInfo 对应表 area_info
type AreaInfo struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Charid int    `gorm:"column:charid" json:"charid,omitempty"`
	Area   int    `gorm:"column:area" json:"area,omitempty"`
	Info   string `gorm:"column:info" json:"info,omitempty"`
}

func (AreaInfo) TableName() string {
	return "area_info"
}
