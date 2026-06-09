package model

import (
	"time"

	"gorm.io/gorm"
)

// Marriage 对应表 marriages
type Marriage struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Marriageid int64          `gorm:"column:marriageid" json:"marriageid"`
	Husbandid  int64          `gorm:"column:husbandid" json:"husbandid"`
	Wifeid     int64          `gorm:"column:wifeid" json:"wifeid"`
}

func (Marriage) TableName() string {
	return "marriages"
}
