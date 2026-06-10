package model

import (
	"time"

	"gorm.io/gorm"
)

// AutobanConfig 对应表 autoban_config
type AutobanConfig struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Type        string    `gorm:"column:type" json:"type,omitempty"`
	Disabled    bool      `gorm:"column:disabled" json:"disabled,omitempty"`
	Points      *int      `gorm:"column:points" json:"points,omitempty"`
	ExpireTime  *int64    `gorm:"column:expireTime" json:"expireTime,omitempty"`
	Description string    `gorm:"column:description" json:"description,omitempty"`
	CreateTime  time.Time `gorm:"column:createTime" json:"createTime,omitempty"`
	UpdateTime  time.Time `gorm:"column:updateTime" json:"updateTime,omitempty"`
}

func (AutobanConfig) TableName() string {
	return "autoban_config"
}
