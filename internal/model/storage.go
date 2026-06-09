package model

import (
	"time"

	"gorm.io/gorm"
)

// Storage mapped from table "storages"
type Storage struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Storageid int64          `gorm:"column:storageid;autoIncrement" json:"storageid,omitempty"`
	Accountid int            `gorm:"column:accountid" json:"accountid,omitempty"`
	World     int            `gorm:"column:world" json:"world,omitempty"`
	Slots     int            `gorm:"column:slots" json:"slots,omitempty"`
	Meso      int            `gorm:"column:meso" json:"meso,omitempty"`
}

func (Storage) TableName() string {
	return "storages"
}
