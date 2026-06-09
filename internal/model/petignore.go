package model

import (
	"time"

	"gorm.io/gorm"
)

// Petignore 映射 petignores 表
type Petignore struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Petid     int            `json:"petid,omitempty"`
	Itemid    int            `json:"itemid,omitempty"`
}

// TableName 指定表名
func (Petignore) TableName() string {
	return "petignores"
}
