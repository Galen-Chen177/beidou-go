package model

import (
	"time"

	"gorm.io/gorm"
)

// Nxcoupon 映射 nxcoupons 表
type Nxcoupon struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Couponid  int            `json:"couponid,omitempty"`
	Rate      int            `json:"rate,omitempty"`
	Activeday int            `json:"activeday,omitempty"`
	Starthour int            `json:"starthour,omitempty"`
	Endhour   int            `json:"endhour,omitempty"`
}

// TableName 指定表名
func (Nxcoupon) TableName() string {
	return "nxcoupons"
}
