package model

import (
	"time"

	"gorm.io/gorm"
)

// BosslogWeekly 对应表 bosslog_weekly
type BosslogWeekly struct {
	gorm.Model

	Characterid int       `gorm:"column:characterid" json:"characterid,omitempty"`
	Bosstype    string    `gorm:"column:bosstype;type:varchar(200)" json:"bosstype,omitempty"`
	Attempttime time.Time `gorm:"column:attempttime" json:"attempttime,omitempty"`
}

func (BosslogWeekly) TableName() string {
	return "bosslog_weekly"
}
