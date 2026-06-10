package model

import (
	"time"

	"gorm.io/gorm"
)

// Savedlocation mapped from table "savedlocations"
type Savedlocation struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	LocID        int            `gorm:"column:id;autoIncrement" json:"locID,omitempty"`
	Characterid  int            `gorm:"column:characterid" json:"characterid,omitempty"`
	Locationtype string         `gorm:"column:locationtype" json:"locationtype,omitempty"`
	Map          int            `gorm:"column:map" json:"map,omitempty"`
	Portal       int            `gorm:"column:portal" json:"portal,omitempty"`
}

func (Savedlocation) TableName() string {
	return "savedlocations"
}
