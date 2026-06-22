package model

import "gorm.io/gorm"

// Inventoryitem 实体映射
type Inventoryitem struct {
	gorm.Model
	Inventoryitemid int64  `gorm:"column:inventoryitemid" json:"inventoryitemid,omitempty"`
	Type            *int   `json:"type,omitempty"`
	Characterid     *int   `json:"characterid,omitempty"`
	Accountid       *int   `json:"accountid,omitempty"`
	Itemid          *int   `json:"itemid,omitempty"`
	Inventorytype   *int   `json:"inventorytype,omitempty"`
	Position        *int   `json:"position,omitempty"`
	Quantity        *int   `json:"quantity,omitempty"`
	Owner           string `gorm:"column:owner;type:varchar(200)" json:"owner,omitempty"`
	Petid           *int   `json:"petid,omitempty"`
	Flag            *int   `json:"flag,omitempty"`
	Expiration      int64  `json:"expiration,omitempty"`
	GiftFrom        string `gorm:"column:giftFrom;type:varchar(200)" json:"giftFrom,omitempty"`
}

func (Inventoryitem) TableName() string {
	return "inventoryitems"
}
