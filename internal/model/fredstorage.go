package model

import (
	"time"

	"gorm.io/gorm"
)

// Fredstorage 实体映射
type Fredstorage struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Cid       int64          `json:"cid,omitempty"`
	Daynotes  int64          `json:"daynotes,omitempty"`
	Timestamp time.Time      `json:"timestamp,omitempty"`
}

func (Fredstorage) TableName() string {
	return "fredstorage"
}
