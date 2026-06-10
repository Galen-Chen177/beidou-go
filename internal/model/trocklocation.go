package model

import (
	"time"

	"gorm.io/gorm"
)

// Trocklocation mapped from table "trocklocations"
type Trocklocation struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Trockid     int            `gorm:"column:trockid;autoIncrement" json:"trockid,omitempty"`
	Characterid int            `gorm:"column:characterid" json:"characterid,omitempty"`
	Mapid       int            `gorm:"column:mapid" json:"mapid,omitempty"`
	Vip         int            `gorm:"column:vip" json:"vip,omitempty"`
}

func (Trocklocation) TableName() string {
	return "trocklocations"
}
