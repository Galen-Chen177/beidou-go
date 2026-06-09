package model

import (
	"time"

	"gorm.io/gorm"
)

// Monsterbook 对应表 monsterbook (复合主键：charid + cardid)
type Monsterbook struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Charid    int            `gorm:"column:charid" json:"charid"`
	Cardid    int            `gorm:"column:cardid" json:"cardid"`
	Level     int            `gorm:"column:level" json:"level"`
}

func (Monsterbook) TableName() string {
	return "monsterbook"
}
