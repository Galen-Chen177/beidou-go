package model

import (
	"time"

	"gorm.io/gorm"
)

// Newyear 映射 newyear 表
type Newyear struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Senderid       int            `json:"senderid,omitempty"`
	Sendername     string         `json:"sendername,omitempty"`
	Receiverid     int            `json:"receiverid,omitempty"`
	Receivername   string         `json:"receivername,omitempty"`
	Message        string         `json:"message,omitempty"`
	Senderdiscard  bool           `json:"senderdiscard,omitempty"`
	Receiverdiscard bool          `json:"receiverdiscard,omitempty"`
	Received       bool           `json:"received,omitempty"`
	Timesent       int64          `json:"timesent,omitempty"`
	Timereceived   int64          `json:"timereceived,omitempty"`
}

// TableName 指定表名
func (Newyear) TableName() string {
	return "newyear"
}
