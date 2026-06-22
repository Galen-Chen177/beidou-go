package model

import (
	"gorm.io/gorm"
)

// Buddy 对应表 buddies
type Buddy struct {
	gorm.Model

	Characterid int    `gorm:"column:characterid" json:"characterid,omitempty"`
	Buddyid     int    `gorm:"column:buddyid" json:"buddyid,omitempty"`
	Pending     int    `gorm:"column:pending" json:"pending,omitempty"`
	Group       string `gorm:"column:group;type:varchar(200)" json:"group,omitempty"`
}

func (Buddy) TableName() string {
	return "buddies"
}
