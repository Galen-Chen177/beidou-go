package model

import (
	"time"

	"gorm.io/gorm"
)

// Questaction 映射 questactions 表
type Questaction struct {
	Questactionid int64          `gorm:"primaryKey" json:"questactionid"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Questid       int            `json:"questid,omitempty"`
	Status        int            `json:"status,omitempty"`
	Data          []byte         `json:"data,omitempty"`
}

// TableName 指定表名
func (Questaction) TableName() string {
	return "questactions"
}
