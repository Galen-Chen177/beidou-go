package model

import (
	"time"

	"gorm.io/gorm"
)

// Makerreagentdata 对应表 makerreagentdata
type Makerreagentdata struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Itemid    int            `gorm:"column:itemid" json:"itemid"`
	Stat      string         `gorm:"column:stat" json:"stat"`
	Value     int            `gorm:"column:value" json:"value"`
}

func (Makerreagentdata) TableName() string {
	return "makerreagentdata"
}
