package model

import (
	"time"

	"gorm.io/gorm"
)

// BbsReply 对应表 bbs_replies
type BbsReply struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Replyid   int64  `gorm:"primaryKey;column:replyid" json:"replyid,omitempty"`
	Threadid  int64  `gorm:"column:threadid" json:"threadid,omitempty"`
	Postercid int64  `gorm:"column:postercid" json:"postercid,omitempty"`
	Timestamp string `gorm:"column:timestamp" json:"timestamp,omitempty"`
	Content   string `gorm:"column:content" json:"content,omitempty"`
}

func (BbsReply) TableName() string {
	return "bbs_replies"
}
