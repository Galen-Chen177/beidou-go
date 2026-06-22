package model

import (
	"gorm.io/gorm"
)

// Alliance 对应表 alliance
type Alliance struct {
	gorm.Model

	Name     string `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	Capacity int64  `gorm:"column:capacity" json:"capacity,omitempty"`
	Notice   string `gorm:"column:notice;type:varchar(200)" json:"notice,omitempty"`
	Rank1    string `gorm:"column:rank1;type:varchar(200)" json:"rank1,omitempty"`
	Rank2    string `gorm:"column:rank2;type:varchar(200)" json:"rank2,omitempty"`
	Rank3    string `gorm:"column:rank3;type:varchar(200)" json:"rank3,omitempty"`
	Rank4    string `gorm:"column:rank4;type:varchar(200)" json:"rank4,omitempty"`
	Rank5    string `gorm:"column:rank5;type:varchar(200)" json:"rank5,omitempty"`
}

func (Alliance) TableName() string {
	return "alliance"
}
