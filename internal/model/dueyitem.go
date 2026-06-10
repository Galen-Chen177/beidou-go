package model

import (
	"time"

	"gorm.io/gorm"
)

// Dueyitem 对应表 dueyitems
type Dueyitem struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Packageid       int64 `gorm:"column:packageid" json:"packageid,omitempty"`
	Inventoryitemid int64 `gorm:"column:inventoryitemid" json:"inventoryitemid,omitempty"`
}

func (Dueyitem) TableName() string {
	return "dueyitems"
}
