package model

import (
	"time"

	"gorm.io/gorm"
)

// Questrequirement 映射 questrequirements 表
type Questrequirement struct {
	Questrequirementid int64          `gorm:"primaryKey" json:"questrequirementid"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Questid            int            `json:"questid,omitempty"`
	Status             int            `json:"status,omitempty"`
	Data               []byte         `json:"data,omitempty"`
}

// TableName 指定表名
func (Questrequirement) TableName() string {
	return "questrequirements"
}
