package model

import (
	"gorm.io/gorm"
)

// BbsThread 对应表 bbs_threads
type BbsThread struct {
	gorm.Model

	Threadid      int64  `gorm:"primaryKey;column:threadid" json:"threadid,omitempty"`
	Postercid     int64  `gorm:"column:postercid" json:"postercid,omitempty"`
	Name          string `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	Timestamp     string `gorm:"column:timestamp;type:varchar(200)" json:"timestamp,omitempty"`
	Icon          int    `gorm:"column:icon" json:"icon,omitempty"`
	Replycount    int    `gorm:"column:replycount" json:"replycount,omitempty"`
	Startpost     string `gorm:"column:startpost;type:varchar(200)" json:"startpost,omitempty"`
	Guildid       int64  `gorm:"column:guildid" json:"guildid,omitempty"`
	Localthreadid int64  `gorm:"column:localthreadid" json:"localthreadid,omitempty"`
}

func (BbsThread) TableName() string {
	return "bbs_threads"
}
