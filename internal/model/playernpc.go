package model

import (
	"time"

	"gorm.io/gorm"
)

// Playernpc 映射 playernpcs 表
type Playernpc struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name         string         `json:"name,omitempty"`
	Hair         int            `json:"hair,omitempty"`
	Face         int            `json:"face,omitempty"`
	Skin         int            `json:"skin,omitempty"`
	Gender       int            `json:"gender,omitempty"`
	X            int            `json:"x,omitempty"`
	Cy           int            `json:"cy,omitempty"`
	World        int            `json:"world,omitempty"`
	Map          int            `json:"map,omitempty"`
	Dir          int            `json:"dir,omitempty"`
	Scriptid     int            `json:"scriptid,omitempty"`
	Fh           int            `json:"fh,omitempty"`
	Rx0          int            `json:"rx0,omitempty"`
	Rx1          int            `json:"rx1,omitempty"`
	Worldrank    int            `json:"worldrank,omitempty"`
	Overallrank  int            `json:"overallrank,omitempty"`
	Worldjobrank int            `json:"worldjobrank,omitempty"`
	Job          int            `json:"job,omitempty"`
}

// TableName 指定表名
func (Playernpc) TableName() string {
	return "playernpcs"
}
