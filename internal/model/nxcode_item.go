package model

import (
	"gorm.io/gorm"
)

// NxcodeItem 映射 nxcode_items 表
type NxcodeItem struct {
	gorm.Model
	Codeid   int `json:"codeid,omitempty"`
	Type     int `json:"type,omitempty"`
	Item     int `json:"item,omitempty"`
	Quantity int `json:"quantity,omitempty"`
}

// TableName 指定表名
func (NxcodeItem) TableName() string {
	return "nxcode_items"
}
