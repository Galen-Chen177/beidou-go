package model

import (
	"time"

	"gorm.io/gorm"
)

// Queststatus mapped from table "queststatus"
type Queststatus struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Queststatusid int64         `gorm:"column:queststatusid;autoIncrement" json:"queststatusid,omitempty"`
	Characterid int            `gorm:"column:characterid" json:"characterid,omitempty"`
	Quest       int            `gorm:"column:quest" json:"quest,omitempty"`
	Status      int            `gorm:"column:status" json:"status,omitempty"`
	Time        int            `gorm:"column:time" json:"time,omitempty"`
	Expires     int64          `gorm:"column:expires" json:"expires,omitempty"`
	Forfeited   int            `gorm:"column:forfeited" json:"forfeited,omitempty"`
	Completed   int            `gorm:"column:completed" json:"completed,omitempty"`
	Info        int            `gorm:"column:info" json:"info,omitempty"`
}

func (Queststatus) TableName() string {
	return "queststatus"
}
