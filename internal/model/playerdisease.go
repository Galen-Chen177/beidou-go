package model

import (
	"gorm.io/gorm"
)

// Playerdisease 映射 playerdiseases 表
type Playerdisease struct {
	gorm.Model
	Charid     int   `json:"charid,omitempty"`
	Disease    int   `json:"disease,omitempty"`
	Mobskillid int   `json:"mobskillid,omitempty"`
	Mobskilllv int   `json:"mobskilllv,omitempty"`
	Length     int64 `json:"length,omitempty"`
}

// TableName 指定表名
func (Playerdisease) TableName() string {
	return "playerdiseases"
}
