package model

import (
	"gorm.io/gorm"
)

// Quickslotkeymapped mapped from table "quickslotkeymapped"
type Quickslotkeymapped struct {
	gorm.Model
	Accountid int   `gorm:"column:accountid" json:"accountid,omitempty"`
	Keymap    int64 `gorm:"column:keymap" json:"keymap,omitempty"`
}

func (Quickslotkeymapped) TableName() string {
	return "quickslotkeymapped"
}
