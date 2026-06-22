package model

import (
	"gorm.io/gorm"
)

// Plife 映射 plife 表
type Plife struct {
	gorm.Model
	World   int    `json:"world,omitempty"`
	Map     int    `json:"map,omitempty"`
	Life    int    `json:"life,omitempty"`
	Type    string `gorm:"column:type;type:varchar(200)" json:"type,omitempty"`
	Cy      int    `json:"cy,omitempty"`
	F       int    `json:"f,omitempty"`
	Fh      int    `json:"fh,omitempty"`
	Rx0     int    `json:"rx0,omitempty"`
	Rx1     int    `json:"rx1,omitempty"`
	X       int    `json:"x,omitempty"`
	Y       int    `json:"y,omitempty"`
	Hide    int    `json:"hide,omitempty"`
	Mobtime int    `json:"mobtime,omitempty"`
	Team    int    `json:"team,omitempty"`
}

// TableName 指定表名
func (Plife) TableName() string {
	return "plife"
}
