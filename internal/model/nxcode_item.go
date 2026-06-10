package model

import (
	"time"

	"gorm.io/gorm"
)

// NxcodeItem 映射 nxcode_items 表
type NxcodeItem struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Codeid    int            `json:"codeid,omitempty"`
	Type      int            `json:"type,omitempty"`
	Item      int            `json:"item,omitempty"`
	Quantity  int            `json:"quantity,omitempty"`
}

// TableName 指定表名
func (NxcodeItem) TableName() string {
	return "nxcode_items"
}
