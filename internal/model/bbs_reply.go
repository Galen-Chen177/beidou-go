package model

import (
	"gorm.io/gorm"
)

// BbsReply 对应表 bbs_replies
type BbsReply struct {
	gorm.Model

	Replyid   int64  `gorm:"primaryKey;column:replyid" json:"replyid,omitempty"`
	Threadid  int64  `gorm:"column:threadid" json:"threadid,omitempty"`
	Postercid int64  `gorm:"column:postercid" json:"postercid,omitempty"`
	Timestamp string `gorm:"column:timestamp;type:varchar(200)" json:"timestamp,omitempty"`
	Content   string `gorm:"column:content;type:varchar(200)" json:"content,omitempty"`
}

func (BbsReply) TableName() string {
	return "bbs_replies"
}
