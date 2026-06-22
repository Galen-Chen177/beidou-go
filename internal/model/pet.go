package model

import "gorm.io/gorm"

// Pet 映射 pets 表
type Pet struct {
	gorm.Model
	Name      string `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	Level     int64  `json:"level,omitempty"`
	Closeness int64  `json:"closeness,omitempty"`
	Fullness  int64  `json:"fullness,omitempty"`
	Summoned  bool   `json:"summoned,omitempty"`
	Flag      int64  `json:"flag,omitempty"`
}

// TableName 指定表名
func (Pet) TableName() string {
	return "pets"
}
