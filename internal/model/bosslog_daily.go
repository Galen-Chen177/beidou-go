package model

import (
	"time"

	"gorm.io/gorm"
)

// BosslogDaily 对应表 bosslog_daily
type BosslogDaily struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Characterid int       `gorm:"column:characterid" json:"characterid,omitempty"`
	Bosstype    string    `gorm:"column:bosstype" json:"bosstype,omitempty"`
	Attempttime time.Time `gorm:"column:attempttime" json:"attempttime,omitempty"`
}

func (BosslogDaily) TableName() string {
	return "bosslog_daily"
}
