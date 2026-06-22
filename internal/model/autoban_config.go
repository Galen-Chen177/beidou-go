package model

import (
	"time"

	"gorm.io/gorm"
)

// AutobanConfig 对应表 autoban_config
type AutobanConfig struct {
	gorm.Model

	Type        string    `gorm:"column:type;type:varchar(200)" json:"type,omitempty"`
	Disabled    bool      `gorm:"column:disabled" json:"disabled,omitempty"`
	Points      *int      `gorm:"column:points" json:"points,omitempty"`
	ExpireTime  *int64    `gorm:"column:expireTime" json:"expireTime,omitempty"`
	Description string    `gorm:"column:description;type:varchar(200)" json:"description,omitempty"`
	CreateTime  time.Time `gorm:"column:createTime" json:"createTime,omitempty"`
	UpdateTime  time.Time `gorm:"column:updateTime" json:"updateTime,omitempty"`
}

func (AutobanConfig) TableName() string {
	return "autoban_config"
}
