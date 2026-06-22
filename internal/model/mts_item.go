package model

import (
	"gorm.io/gorm"
)

// MtsItem 对应表 mts_items
type MtsItem struct {
	gorm.Model
	Tab          int    `gorm:"column:tab" json:"tab"`
	Type         int    `gorm:"column:type" json:"type"`
	Itemid       int64  `gorm:"column:itemid" json:"itemid"`
	Quantity     int    `gorm:"column:quantity" json:"quantity"`
	Seller       int    `gorm:"column:seller" json:"seller"`
	Price        int    `gorm:"column:price" json:"price"`
	BidIncre     int    `gorm:"column:bidIncre" json:"bidIncre"`
	BuyNow       int    `gorm:"column:buyNow" json:"buyNow"`
	Position     int    `gorm:"column:position" json:"position"`
	Upgradeslots int    `gorm:"column:upgradeslots" json:"upgradeslots"`
	Level        int    `gorm:"column:level" json:"level"`
	Itemlevel    int    `gorm:"column:itemlevel" json:"itemlevel"`
	Itemexp      int64  `gorm:"column:itemexp" json:"itemexp"`
	Ringid       int    `gorm:"column:ringid" json:"ringid"`
	Str          int    `gorm:"column:str" json:"str"`
	Dex          int    `gorm:"column:dex" json:"dex"`
	Inte         int    `gorm:"column:int" json:"inte"`
	Luk          int    `gorm:"column:luk" json:"luk"`
	Hp           int    `gorm:"column:hp" json:"hp"`
	Mp           int    `gorm:"column:mp" json:"mp"`
	Watk         int    `gorm:"column:watk" json:"watk"`
	Matk         int    `gorm:"column:matk" json:"matk"`
	Wdef         int    `gorm:"column:wdef" json:"wdef"`
	Mdef         int    `gorm:"column:mdef" json:"mdef"`
	Acc          int    `gorm:"column:acc" json:"acc"`
	Avoid        int    `gorm:"column:avoid" json:"avoid"`
	Hands        int    `gorm:"column:hands" json:"hands"`
	Speed        int    `gorm:"column:speed" json:"speed"`
	Jump         int    `gorm:"column:jump" json:"jump"`
	Locked       int    `gorm:"column:locked" json:"locked"`
	Isequip      int    `gorm:"column:isequip" json:"isequip"`
	Owner        string `gorm:"column:owner;type:varchar(200)" json:"owner"`
	Sellername   string `gorm:"column:sellername;type:varchar(200)" json:"sellername"`
	SellEnds     string `gorm:"column:sellEnds;type:varchar(200)" json:"sellEnds"`
	Transfer     int    `gorm:"column:transfer" json:"transfer"`
	Vicious      int64  `gorm:"column:vicious" json:"vicious"`
	Flag         int64  `gorm:"column:flag" json:"flag"`
	Expiration   int64  `gorm:"column:expiration" json:"expiration"`
	GiftFrom     string `gorm:"column:giftFrom;type:varchar(200)" json:"giftFrom"`
}

func (MtsItem) TableName() string {
	return "mts_items"
}
