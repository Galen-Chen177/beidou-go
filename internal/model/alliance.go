package model

import (
	"time"

	"gorm.io/gorm"
)

// Alliance 对应表 alliance
type Alliance struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name     string `gorm:"column:name" json:"name,omitempty"`
	Capacity int64  `gorm:"column:capacity" json:"capacity,omitempty"`
	Notice   string `gorm:"column:notice" json:"notice,omitempty"`
	Rank1    string `gorm:"column:rank1" json:"rank1,omitempty"`
	Rank2    string `gorm:"column:rank2" json:"rank2,omitempty"`
	Rank3    string `gorm:"column:rank3" json:"rank3,omitempty"`
	Rank4    string `gorm:"column:rank4" json:"rank4,omitempty"`
	Rank5    string `gorm:"column:rank5" json:"rank5,omitempty"`
}

func (Alliance) TableName() string {
	return "alliance"
}
