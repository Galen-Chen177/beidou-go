package model

import (
	"gorm.io/gorm"
)

// Petignore 映射 petignores 表
type Petignore struct {
	gorm.Model
	Petid  int `json:"petid,omitempty"`
	Itemid int `json:"itemid,omitempty"`
}

// TableName 指定表名
func (Petignore) TableName() string {
	return "petignores"
}
