package model

import (
	"time"

	"gorm.io/gorm"
)

// PlayernpcEquip 映射 playernpcs_equip 表
type PlayernpcEquip struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Npcid     int            `json:"npcid,omitempty"`
	Equipid   int            `json:"equipid,omitempty"`
	Type      int            `json:"type,omitempty"`
	Equippos  int16          `json:"equippos,omitempty"`
}

// TableName 指定表名
func (PlayernpcEquip) TableName() string {
	return "playernpcs_equip"
}
