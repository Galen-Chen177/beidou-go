package model

import (
	"gorm.io/gorm"
)

// Note 映射 notes 表
type Note struct {
	gorm.Model
	To        string `gorm:"column:to;type:varchar(200)" json:"to,omitempty"`
	From      string `gorm:"column:from;type:varchar(200)" json:"from,omitempty"`
	Message   string `gorm:"column:message;type:varchar(200)" json:"message,omitempty"`
	Timestamp int64  `json:"timestamp,omitempty"`
	Fame      int    `json:"fame,omitempty"`
	Deleted   int    `json:"deleted,omitempty"`
}

// TableName 指定表名
func (Note) TableName() string {
	return "notes"
}
