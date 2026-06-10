package model

import (
	"time"

	"gorm.io/gorm"
)

// Plife 映射 plife 表
type Plife struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	World     int            `json:"world,omitempty"`
	Map       int            `json:"map,omitempty"`
	Life      int            `json:"life,omitempty"`
	Type      string         `json:"type,omitempty"`
	Cy        int            `json:"cy,omitempty"`
	F         int            `json:"f,omitempty"`
	Fh        int            `json:"fh,omitempty"`
	Rx0       int            `json:"rx0,omitempty"`
	Rx1       int            `json:"rx1,omitempty"`
	X         int            `json:"x,omitempty"`
	Y         int            `json:"y,omitempty"`
	Hide      int            `json:"hide,omitempty"`
	Mobtime   int            `json:"mobtime,omitempty"`
	Team      int            `json:"team,omitempty"`
}

// TableName 指定表名
func (Plife) TableName() string {
	return "plife"
}
