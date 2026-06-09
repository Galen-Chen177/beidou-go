package model

import (
	"time"

	"gorm.io/gorm"
)

// Reactordrop mapped from table "reactordrops"
type Reactordrop struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Reactordropid int64          `gorm:"column:reactordropid;autoIncrement" json:"reactordropid,omitempty"`
	Reactorid     int            `gorm:"column:reactorid" json:"reactorid,omitempty"`
	Itemid        int            `gorm:"column:itemid" json:"itemid,omitempty"`
	Chance        int            `gorm:"column:chance" json:"chance,omitempty"`
	Questid       int            `gorm:"column:questid" json:"questid,omitempty"`
}

func (Reactordrop) TableName() string {
	return "reactordrops"
}
