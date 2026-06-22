package model

import (
	"gorm.io/gorm"
)

// ModifiedCashItem 对应表 modified_cash_item
type ModifiedCashItem struct {
	gorm.Model
	Sn              int   `gorm:"column:sn" json:"sn"`
	ItemId          int   `gorm:"column:itemId" json:"itemId"`
	Count           int16 `gorm:"column:count" json:"count"`
	Price           int   `gorm:"column:price" json:"price"`
	Bonus           int   `gorm:"column:bonus" json:"bonus"`
	Priority        int   `gorm:"column:priority" json:"priority"`
	Period          int64 `gorm:"column:period" json:"period"`
	MaplePoint      int   `gorm:"column:maplePoint" json:"maplePoint"`
	Meso            int   `gorm:"column:meso" json:"meso"`
	ForPremiumUser  int   `gorm:"column:forPremiumUser" json:"forPremiumUser"`
	CommodityGender int   `gorm:"column:commodityGender" json:"commodityGender"`
	OnSale          int   `gorm:"column:onSale" json:"onSale"`
	Clz             int   `gorm:"column:class" json:"clz"`
	Limit           int   `gorm:"column:limit" json:"limit"`
	PbCash          int   `gorm:"column:pbCash" json:"pbCash"`
	PbPoint         int   `gorm:"column:pbPoint" json:"pbPoint"`
	PbGift          int   `gorm:"column:pbGift" json:"pbGift"`
	PackageSn       int   `gorm:"column:packageSn" json:"packageSn"`
}

func (ModifiedCashItem) TableName() string {
	return "modified_cash_item"
}
