package model

import (
	"time"

	"gorm.io/gorm"
)

// Medalmap 对应表 medalmaps
type Medalmap struct {
	ID            uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt     *time.Time     `json:"createdAt"`
	UpdatedAt     *time.Time     `json:"-"`
	DeletedAt     gorm.DeletedAt `gorm:"index"`
	Characterid   int            `gorm:"column:characterid" json:"characterid"`
	Queststatusid int64          `gorm:"column:queststatusid" json:"queststatusid"`
	Mapid         int            `gorm:"column:mapid" json:"mapid"`
}

func (Medalmap) TableName() string {
	return "medalmaps"
}
