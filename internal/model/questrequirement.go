package model

import (
	"gorm.io/gorm"
)

// Questrequirement 映射 questrequirements 表
type Questrequirement struct {
	gorm.Model
	Questid int    `json:"questid,omitempty"`
	Status  int    `json:"status,omitempty"`
	Data    []byte `json:"data,omitempty"`
}

// TableName 指定表名
func (Questrequirement) TableName() string {
	return "questrequirements"
}
