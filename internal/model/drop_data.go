package model

import (
	"time"

	"gorm.io/gorm"
)

// DropData 对应表 drop_data
type DropData struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Dropperid       int `gorm:"column:dropperid" json:"dropperid,omitempty"`
	Itemid          int `gorm:"column:itemid" json:"itemid,omitempty"`
	MinimumQuantity int `gorm:"column:minimumQuantity" json:"minimumQuantity,omitempty"`
	MaximumQuantity int `gorm:"column:maximumQuantity" json:"maximumQuantity,omitempty"`
	Questid         int `gorm:"column:questid" json:"questid,omitempty"`
	Chance          int `gorm:"column:chance" json:"chance,omitempty"`
}

func (DropData) TableName() string {
	return "drop_data"
}
