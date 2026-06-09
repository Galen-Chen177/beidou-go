package model

import (
	"time"

	"gorm.io/gorm"
)

// GachaponReward 实体映射
type GachaponReward struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	PoolId    *int           `json:"poolId,omitempty"`
	ItemId    *int           `json:"itemId,omitempty"`
	ItemName  string         `gorm:"-" json:"itemName,omitempty"`
	Quantity  int16          `json:"quantity,omitempty"`
	CreateTime time.Time     `gorm:"column:create_time" json:"createTime,omitempty"`
	Comment   string         `json:"comment,omitempty"`
}

func (GachaponReward) TableName() string {
	return "gachapon_reward"
}
