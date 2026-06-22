package model

import (
	"time"

	"gorm.io/gorm"
)

// Fredstorage 实体映射
type Fredstorage struct {
	gorm.Model
	Cid       int64     `json:"cid,omitempty"`
	Daynotes  int64     `json:"daynotes,omitempty"`
	Timestamp time.Time `json:"timestamp,omitempty"`
}

func (Fredstorage) TableName() string {
	return "fredstorage"
}
