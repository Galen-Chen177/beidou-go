package model

import (
	"gorm.io/gorm"
)

// Skillmacro mapped from table "skillmacros"
type Skillmacro struct {
	gorm.Model
	Characterid int    `gorm:"column:characterid" json:"characterid,omitempty"`
	Position    int    `gorm:"column:position" json:"position,omitempty"`
	Skill1      int    `gorm:"column:skill1" json:"skill1,omitempty"`
	Skill2      int    `gorm:"column:skill2" json:"skill2,omitempty"`
	Skill3      int    `gorm:"column:skill3" json:"skill3,omitempty"`
	Name        string `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	Shout       int    `gorm:"column:shout" json:"shout,omitempty"`
}

func (Skillmacro) TableName() string {
	return "skillmacros"
}
