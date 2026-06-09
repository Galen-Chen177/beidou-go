package model

import (
	"time"

	"gorm.io/gorm"
)

// ServerQueue mapped from table "server_queue"
type ServerQueue struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	QueueID     int            `gorm:"column:id;autoIncrement" json:"queueID,omitempty"`
	Accountid   int            `gorm:"column:accountid" json:"accountid,omitempty"`
	Characterid int            `gorm:"column:characterid" json:"characterid,omitempty"`
	Type        int            `gorm:"column:type" json:"type,omitempty"`
	Value       int            `gorm:"column:value" json:"value,omitempty"`
	Message     string         `gorm:"column:message" json:"message,omitempty"`
	CreateTime  time.Time      `gorm:"column:createTime" json:"createTime,omitempty"`
}

func (ServerQueue) TableName() string {
	return "server_queue"
}
