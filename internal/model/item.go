package model

import "time"

// Item 物品模型
type Item struct {
	ID          int64     `gorm:"primaryKey;column:id;autoIncrement"`
	CharacterID int32     `gorm:"column:characterid;index"`
	ItemID      int32     `gorm:"column:itemid"`
	Slot        int16     `gorm:"column:slot"`        // 背包格子位置（-1 表示已装备）
	Quantity    int16     `gorm:"column:quantity"`
	Expiration  time.Time `gorm:"column:expiration"`
	Owner       string    `gorm:"column:owner;size:13"`
}

func (Item) TableName() string {
	return "inventoryitems"
}
