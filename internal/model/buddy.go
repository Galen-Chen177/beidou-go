package model

import (
	"time"

	"gorm.io/gorm"
)

// Buddy 对应表 buddies
type Buddy struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Characterid int    `gorm:"column:characterid" json:"characterid,omitempty"`
	Buddyid     int    `gorm:"column:buddyid" json:"buddyid,omitempty"`
	Pending     int    `gorm:"column:pending" json:"pending,omitempty"`
	Group       string `gorm:"column:group" json:"group,omitempty"`
}

func (Buddy) TableName() string {
	return "buddies"
}
