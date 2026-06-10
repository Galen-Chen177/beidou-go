package model

import (
	"time"

	"gorm.io/gorm"
)

// Skillmacro mapped from table "skillmacros"
type Skillmacro struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	MacroID     int            `gorm:"column:id;autoIncrement" json:"macroID,omitempty"`
	Characterid int            `gorm:"column:characterid" json:"characterid,omitempty"`
	Position    int            `gorm:"column:position" json:"position,omitempty"`
	Skill1      int            `gorm:"column:skill1" json:"skill1,omitempty"`
	Skill2      int            `gorm:"column:skill2" json:"skill2,omitempty"`
	Skill3      int            `gorm:"column:skill3" json:"skill3,omitempty"`
	Name        string         `gorm:"column:name" json:"name,omitempty"`
	Shout       int            `gorm:"column:shout" json:"shout,omitempty"`
}

func (Skillmacro) TableName() string {
	return "skillmacros"
}
