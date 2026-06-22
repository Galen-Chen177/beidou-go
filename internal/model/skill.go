package model

import (
	"gorm.io/gorm"
)

// Skill mapped from table "skills"
type Skill struct {
	gorm.Model
	Skillid     int   `gorm:"column:skillid" json:"skillid,omitempty"`
	Characterid int   `gorm:"column:characterid" json:"characterid,omitempty"`
	Skilllevel  int   `gorm:"column:skilllevel" json:"skilllevel,omitempty"`
	Masterlevel int   `gorm:"column:masterlevel" json:"masterlevel,omitempty"`
	Expiration  int64 `gorm:"column:expiration" json:"expiration,omitempty"`
}

func (Skill) TableName() string {
	return "skills"
}
