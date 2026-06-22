package model

import (
	"gorm.io/gorm"
)

// Questaction 映射 questactions 表
type Questaction struct {
	gorm.Model
	Questid int    `json:"questid,omitempty"`
	Status  int    `json:"status,omitempty"`
	Data    []byte `json:"data,omitempty"`
}

// TableName 指定表名
func (Questaction) TableName() string {
	return "questactions"
}
