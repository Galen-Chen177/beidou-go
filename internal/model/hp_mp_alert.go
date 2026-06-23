package model

import "gorm.io/gorm"

// HpMpAlert 实体映射
type HpMpAlert struct {
	gorm.Model
	CID *int `gorm:"column:cID" json:"cID,omitempty"`
	Hp  int8 `json:"hp,omitempty"`
	Mp  int8 `json:"mp,omitempty"`
}

func (HpMpAlert) TableName() string {
	return "hp_mp_alert"
}
