package model

import (
	"time"

	"gorm.io/gorm"
)

// Pet 映射 pets 表
type Pet struct {
	Petid     int64          `gorm:"primaryKey" json:"petid"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name      string         `json:"name,omitempty"`
	Level     int64          `json:"level,omitempty"`
	Closeness int64          `json:"closeness,omitempty"`
	Fullness  int64          `json:"fullness,omitempty"`
	Summoned  bool           `json:"summoned,omitempty"`
	Flag      int64          `json:"flag,omitempty"`
}

// TableName 指定表名
func (Pet) TableName() string {
	return "pets"
}
