package model

import (
	"time"

	"gorm.io/gorm"
)

// Inventorymerchant 对应表 inventorymerchant
type Inventorymerchant struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Inventorymerchantid int64        `gorm:"column:inventorymerchantid" json:"inventorymerchantid"`
	Inventoryitemid   int64          `gorm:"column:inventoryitemid" json:"inventoryitemid"`
	Characterid       int            `gorm:"column:characterid" json:"characterid"`
	Bundles           int            `gorm:"column:bundles" json:"bundles"`
}

func (Inventorymerchant) TableName() string {
	return "inventorymerchant"
}
