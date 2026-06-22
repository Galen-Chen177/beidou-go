package model

import (
	"gorm.io/gorm"
)

// Queststatus mapped from table "queststatus"
type Queststatus struct {
	gorm.Model
	Characterid int   `gorm:"column:characterid" json:"characterid,omitempty"`
	Quest       int   `gorm:"column:quest" json:"quest,omitempty"`
	Status      int   `gorm:"column:status" json:"status,omitempty"`
	Time        int   `gorm:"column:time" json:"time,omitempty"`
	Expires     int64 `gorm:"column:expires" json:"expires,omitempty"`
	Forfeited   int   `gorm:"column:forfeited" json:"forfeited,omitempty"`
	Completed   int   `gorm:"column:completed" json:"completed,omitempty"`
	Info        int   `gorm:"column:info" json:"info,omitempty"`
}

func (Queststatus) TableName() string {
	return "queststatus"
}
