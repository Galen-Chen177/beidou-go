package model

import (
	"time"

	"gorm.io/gorm"
)

// BbsThread 对应表 bbs_threads
type BbsThread struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Threadid      int64  `gorm:"primaryKey;column:threadid" json:"threadid,omitempty"`
	Postercid     int64  `gorm:"column:postercid" json:"postercid,omitempty"`
	Name          string `gorm:"column:name" json:"name,omitempty"`
	Timestamp     string `gorm:"column:timestamp" json:"timestamp,omitempty"`
	Icon          int    `gorm:"column:icon" json:"icon,omitempty"`
	Replycount    int    `gorm:"column:replycount" json:"replycount,omitempty"`
	Startpost     string `gorm:"column:startpost" json:"startpost,omitempty"`
	Guildid       int64  `gorm:"column:guildid" json:"guildid,omitempty"`
	Localthreadid int64  `gorm:"column:localthreadid" json:"localthreadid,omitempty"`
}

func (BbsThread) TableName() string {
	return "bbs_threads"
}
