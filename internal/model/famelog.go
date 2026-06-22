package model

import (
	"time"

	"gorm.io/gorm"
)

// Famelog 实体映射
type Famelog struct {
	gorm.Model
	Famelogid     int       `gorm:"column:famelogid" json:"famelogid,omitempty"`
	Characterid   *int      `json:"characterid,omitempty"`
	CharacteridTo *int      `json:"characteridTo,omitempty"`
	When          time.Time `json:"when,omitempty"`
}

func (Famelog) TableName() string {
	return "famelog"
}
