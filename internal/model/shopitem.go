package model

import (
	"time"

	"gorm.io/gorm"
)

// Shopitem mapped from table "shopitems"
type Shopitem struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Shopitemid int64          `gorm:"column:shopitemid;autoIncrement" json:"shopitemid,omitempty"`
	Shopid     int64          `gorm:"column:shopid" json:"shopid,omitempty"`
	Itemid     int            `gorm:"column:itemid" json:"itemid,omitempty"`
	Price      int            `gorm:"column:price" json:"price,omitempty"`
	Pitch      int            `gorm:"column:pitch" json:"pitch,omitempty"`
	Position   int            `gorm:"column:position" json:"position,omitempty"`
}

func (Shopitem) TableName() string {
	return "shopitems"
}
