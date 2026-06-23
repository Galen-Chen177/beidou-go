package model

import (
	"time"

	"gorm.io/gorm"
)

// GachaponRewardPool 实体映射
type GachaponRewardPool struct {
	gorm.Model
	Name         string    `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	GachaponID   *int      `json:"gachaponID,omitempty"`
	GachaponName string    `gorm:"-" json:"gachaponName,omitempty"`
	Weight       *int      `json:"weight,omitempty"`
	IsPublic     *bool     `json:"isPublic,omitempty"`
	Prob         *int      `json:"prob,omitempty"`
	StartTime    time.Time `json:"startTime,omitempty"`
	EndTime      time.Time `json:"endTime,omitempty"`
	Notification *bool     `json:"notification,omitempty"`
	Comment      string    `gorm:"column:comment;type:varchar(200)" json:"comment,omitempty"`
}

func (GachaponRewardPool) TableName() string {
	return "gachapon_reward_pool"
}
