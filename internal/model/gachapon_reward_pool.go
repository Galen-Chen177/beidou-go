package model

import (
	"time"

	"gorm.io/gorm"
)

// GachaponRewardPool 实体映射
type GachaponRewardPool struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Name         string         `json:"name,omitempty"`
	GachaponId   *int           `json:"gachaponId,omitempty"`
	GachaponName string         `gorm:"-" json:"gachaponName,omitempty"`
	Weight       *int           `json:"weight,omitempty"`
	IsPublic     *bool          `json:"isPublic,omitempty"`
	Prob         *int           `json:"prob,omitempty"`
	StartTime    time.Time      `json:"startTime,omitempty"`
	EndTime      time.Time      `json:"endTime,omitempty"`
	Notification *bool          `json:"notification,omitempty"`
	Comment      string         `json:"comment,omitempty"`
}

func (GachaponRewardPool) TableName() string {
	return "gachapon_reward_pool"
}
