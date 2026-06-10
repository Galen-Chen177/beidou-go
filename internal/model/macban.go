package model

import (
	"time"

	"gorm.io/gorm"
)

// Macban 对应表 macbans
type Macban struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Macbanid  int64          `gorm:"column:macbanid" json:"macbanid"`
	Mac       string         `gorm:"column:mac" json:"mac"`
	Aid       string         `gorm:"column:aid" json:"aid"`
}

func (Macban) TableName() string {
	return "macbans"
}
