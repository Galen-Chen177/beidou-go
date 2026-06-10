package model

import (
	"time"

	"gorm.io/gorm"
)

// Macfilter 对应表 macfilters
type Macfilter struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Macfilterid int64          `gorm:"column:macfilterid" json:"macfilterid"`
	Filter      string         `gorm:"column:filter" json:"filter"`
}

func (Macfilter) TableName() string {
	return "macfilters"
}
