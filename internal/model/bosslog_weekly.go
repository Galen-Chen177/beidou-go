package model

import (
	"time"

	"gorm.io/gorm"
)

// BosslogWeekly 对应表 bosslog_weekly
type BosslogWeekly struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Characterid int       `gorm:"column:characterid" json:"characterid,omitempty"`
	Bosstype    string    `gorm:"column:bosstype" json:"bosstype,omitempty"`
	Attempttime time.Time `gorm:"column:attempttime" json:"attempttime,omitempty"`
}

func (BosslogWeekly) TableName() string {
	return "bosslog_weekly"
}
