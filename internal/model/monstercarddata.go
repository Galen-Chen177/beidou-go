package model

import (
	"time"

	"gorm.io/gorm"
)

// Monstercarddata 对应表 monstercarddata
type Monstercarddata struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Cardid    int            `gorm:"column:cardid" json:"cardid"`
	Mobid     int            `gorm:"column:mobid" json:"mobid"`
}

func (Monstercarddata) TableName() string {
	return "monstercarddata"
}
