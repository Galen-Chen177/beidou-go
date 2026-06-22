package model

import (
	"gorm.io/gorm"
)

// PlayernpcField 映射 playernpcs_field 表
type PlayernpcField struct {
	gorm.Model
	World  int `json:"world,omitempty"`
	Map    int `json:"map,omitempty"`
	Step   int `json:"step,omitempty"`
	Podium int `json:"podium,omitempty"`
}

// TableName 指定表名
func (PlayernpcField) TableName() string {
	return "playernpcs_field"
}
