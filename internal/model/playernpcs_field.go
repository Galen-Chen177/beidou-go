package model

import (
	"time"

	"gorm.io/gorm"
)

// PlayernpcField 映射 playernpcs_field 表
type PlayernpcField struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	World     int            `json:"world,omitempty"`
	Map       int            `json:"map,omitempty"`
	Step      int            `json:"step,omitempty"`
	Podium    int            `json:"podium,omitempty"`
}

// TableName 指定表名
func (PlayernpcField) TableName() string {
	return "playernpcs_field"
}
