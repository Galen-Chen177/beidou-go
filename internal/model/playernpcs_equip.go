package model

import (
	"gorm.io/gorm"
)

// PlayernpcEquip 映射 playernpcs_equip 表
type PlayernpcEquip struct {
	gorm.Model
	Npcid    int   `json:"npcid,omitempty"`
	Equipid  int   `json:"equipid,omitempty"`
	Type     int   `json:"type,omitempty"`
	Equippos int16 `json:"equippos,omitempty"`
}

// TableName 指定表名
func (PlayernpcEquip) TableName() string {
	return "playernpcs_equip"
}
