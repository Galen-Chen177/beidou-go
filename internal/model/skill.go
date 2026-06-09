package model

import (
	"time"

	"gorm.io/gorm"
)

// Skill mapped from table "skills"
type Skill struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	SkillPK     int            `gorm:"column:id;autoIncrement" json:"skillPK,omitempty"`
	Skillid     int            `gorm:"column:skillid" json:"skillid,omitempty"`
	Characterid int            `gorm:"column:characterid" json:"characterid,omitempty"`
	Skilllevel  int            `gorm:"column:skilllevel" json:"skilllevel,omitempty"`
	Masterlevel int            `gorm:"column:masterlevel" json:"masterlevel,omitempty"`
	Expiration  int64          `gorm:"column:expiration" json:"expiration,omitempty"`
}

func (Skill) TableName() string {
	return "skills"
}
