package model

import (
	"time"

	"gorm.io/gorm"
)

// ServerQueue mapped from table "server_queue"
type ServerQueue struct {
	gorm.Model
	Accountid   int       `gorm:"column:accountid" json:"accountid,omitempty"`
	Characterid int       `gorm:"column:characterid" json:"characterid,omitempty"`
	Type        int       `gorm:"column:type" json:"type,omitempty"`
	Value       int       `gorm:"column:value" json:"value,omitempty"`
	Message     string    `gorm:"column:message;type:varchar(200)" json:"message,omitempty"`
	CreateTime  time.Time `gorm:"column:createTime" json:"createTime,omitempty"`
}

func (ServerQueue) TableName() string {
	return "server_queue"
}
