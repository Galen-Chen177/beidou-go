package model

import (
	"gorm.io/gorm"
)

// Marriage 对应表 marriages
type Marriage struct {
	gorm.Model
	Marriageid int64 `gorm:"column:marriageid" json:"marriageid"`
	Husbandid  int64 `gorm:"column:husbandid" json:"husbandid"`
	Wifeid     int64 `gorm:"column:wifeid" json:"wifeid"`
}

func (Marriage) TableName() string {
	return "marriages"
}
