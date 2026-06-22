package model

import (
	"gorm.io/gorm"
)

// Playernpc 映射 playernpcs 表
type Playernpc struct {
	gorm.Model
	Name         string `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	Hair         int    `json:"hair,omitempty"`
	Face         int    `json:"face,omitempty"`
	Skin         int    `json:"skin,omitempty"`
	Gender       int    `json:"gender,omitempty"`
	X            int    `json:"x,omitempty"`
	Cy           int    `json:"cy,omitempty"`
	World        int    `json:"world,omitempty"`
	Map          int    `json:"map,omitempty"`
	Dir          int    `json:"dir,omitempty"`
	Scriptid     int    `json:"scriptid,omitempty"`
	Fh           int    `json:"fh,omitempty"`
	Rx0          int    `json:"rx0,omitempty"`
	Rx1          int    `json:"rx1,omitempty"`
	Worldrank    int    `json:"worldrank,omitempty"`
	Overallrank  int    `json:"overallrank,omitempty"`
	Worldjobrank int    `json:"worldjobrank,omitempty"`
	Job          int    `json:"job,omitempty"`
}

// TableName 指定表名
func (Playernpc) TableName() string {
	return "playernpcs"
}
