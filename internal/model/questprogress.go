package model

import (
	"gorm.io/gorm"
)

// Questprogress 映射 questprogress 表
type Questprogress struct {
	gorm.Model
	Characterid   int    `json:"characterid,omitempty"`
	Queststatusid int64  `json:"queststatusid,omitempty"`
	Progressid    int    `json:"progressid,omitempty"`
	Progress      string `gorm:"column:progress;type:varchar(200)" json:"progress,omitempty"`
}

// TableName 指定表名
func (Questprogress) TableName() string {
	return "questprogress"
}
