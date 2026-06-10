package model

import (
	"time"

	"gorm.io/gorm"
)

// Ipban 对应表 ipbans
type Ipban struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Ipbanid   int64          `gorm:"column:ipbanid" json:"ipbanid"`
	Ip        string         `gorm:"column:ip" json:"ip"`
	Aid       string         `gorm:"column:aid" json:"aid"`
}

func (Ipban) TableName() string {
	return "ipbans"
}
