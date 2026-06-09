package model

import (
	"time"

	"gorm.io/gorm"
)

// Note 映射 notes 表
type Note struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	To        string         `json:"to,omitempty"`
	From      string         `json:"from,omitempty"`
	Message   string         `json:"message,omitempty"`
	Timestamp int64          `json:"timestamp,omitempty"`
	Fame      int            `json:"fame,omitempty"`
	Deleted   int            `json:"deleted,omitempty"`
}

// TableName 指定表名
func (Note) TableName() string {
	return "notes"
}
