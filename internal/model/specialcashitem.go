package model

import (
	"time"

	"gorm.io/gorm"
)

// Specialcashitem mapped from table "specialcashitems"
type Specialcashitem struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	ItemID    int            `gorm:"column:id" json:"itemID,omitempty"`
	Sn        int            `gorm:"column:sn" json:"sn,omitempty"`
	Modifier  int            `gorm:"column:modifier" json:"modifier,omitempty"`
	Info      int            `gorm:"column:info" json:"info,omitempty"`
}

func (Specialcashitem) TableName() string {
	return "specialcashitems"
}
