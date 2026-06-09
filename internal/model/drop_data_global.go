package model

import (
	"time"

	"gorm.io/gorm"
)

// DropDataGlobal 对应表 drop_data_global
type DropDataGlobal struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Continent       int    `gorm:"column:continent" json:"continent,omitempty"`
	Itemid          int    `gorm:"column:itemid" json:"itemid,omitempty"`
	MinimumQuantity int    `gorm:"column:minimumQuantity" json:"minimumQuantity,omitempty"`
	MaximumQuantity int    `gorm:"column:maximumQuantity" json:"maximumQuantity,omitempty"`
	Questid         int    `gorm:"column:questid" json:"questid,omitempty"`
	Chance          int    `gorm:"column:chance" json:"chance,omitempty"`
	Comments        string `gorm:"column:comments" json:"comments,omitempty"`
}

func (DropDataGlobal) TableName() string {
	return "drop_data_global"
}
