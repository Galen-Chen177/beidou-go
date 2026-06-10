package model

import (
	"time"

	"gorm.io/gorm"
)

// Questprogress 映射 questprogress 表
type Questprogress struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Characterid   int            `json:"characterid,omitempty"`
	Queststatusid int64          `json:"queststatusid,omitempty"`
	Progressid    int            `json:"progressid,omitempty"`
	Progress      string         `json:"progress,omitempty"`
}

// TableName 指定表名
func (Questprogress) TableName() string {
	return "questprogress"
}
