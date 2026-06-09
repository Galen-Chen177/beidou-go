package model

import (
	"time"

	"gorm.io/gorm"
)

// Quickslotkeymapped mapped from table "quickslotkeymapped"
type Quickslotkeymapped struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Accountid int            `gorm:"column:accountid" json:"accountid,omitempty"`
	Keymap    int64          `gorm:"column:keymap" json:"keymap,omitempty"`
}

func (Quickslotkeymapped) TableName() string {
	return "quickslotkeymapped"
}
