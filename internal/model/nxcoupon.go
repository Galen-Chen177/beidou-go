package model

import (
	"gorm.io/gorm"
)

// Nxcoupon 映射 nxcoupons 表
type Nxcoupon struct {
	gorm.Model
	Couponid  int `json:"couponid,omitempty"`
	Rate      int `json:"rate,omitempty"`
	Activeday int `json:"activeday,omitempty"`
	Starthour int `json:"starthour,omitempty"`
	Endhour   int `json:"endhour,omitempty"`
}

// TableName 指定表名
func (Nxcoupon) TableName() string {
	return "nxcoupons"
}
