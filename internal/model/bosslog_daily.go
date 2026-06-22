package model

import (
	"time"

	"gorm.io/gorm"
)

// BosslogDaily 对应表 bosslog_daily
type BosslogDaily struct {
	gorm.Model

	Characterid int       `gorm:"column:characterid" json:"characterid,omitempty"`
	Bosstype    string    `gorm:"column:bosstype;type:varchar(200)" json:"bosstype,omitempty"`
	Attempttime time.Time `gorm:"column:attempttime" json:"attempttime,omitempty"`
}

func (BosslogDaily) TableName() string {
	return "bosslog_daily"
}
