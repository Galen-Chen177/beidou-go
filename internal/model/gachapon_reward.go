package model

import (
	"time"

	"gorm.io/gorm"
)

// GachaponReward 实体映射
type GachaponReward struct {
	gorm.Model
	PoolId     *int      `json:"poolId,omitempty"`
	ItemId     *int      `json:"itemId,omitempty"`
	ItemName   string    `gorm:"-" json:"itemName,omitempty"`
	Quantity   int16     `json:"quantity,omitempty"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime,omitempty"`
	Comment    string    `gorm:"column:comment;type:varchar(200)" json:"comment,omitempty"`
}

func (GachaponReward) TableName() string {
	return "gachapon_reward"
}
