package model

import (
	"gorm.io/gorm"
)

// Nxcode 映射 nxcode 表
type Nxcode struct {
	gorm.Model
	Code       string `gorm:"column:code;type:varchar(200)" json:"code,omitempty"`
	Retriever  string `gorm:"column:retriever;type:varchar(200)" json:"retriever,omitempty"`
	Expiration uint64 `gorm:"column:expiration;type:bigint unsigned;not null;default:0" json:"expiration,omitempty"`
}

// TableName 指定表名
func (Nxcode) TableName() string {
	return "nxcode"
}
