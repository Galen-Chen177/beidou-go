package model

// Inventoryitem 实体映射
type Inventoryitem struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	Inventoryitemid int64  `gorm:"column:inventoryitemid" json:"inventoryitemid,omitempty"`
	Type            *int   `json:"type,omitempty"`
	Characterid     *int   `json:"characterid,omitempty"`
	Accountid       *int   `json:"accountid,omitempty"`
	Itemid          *int   `json:"itemid,omitempty"`
	Inventorytype   *int   `json:"inventorytype,omitempty"`
	Position        *int   `json:"position,omitempty"`
	Quantity        *int   `json:"quantity,omitempty"`
	Owner           string `json:"owner,omitempty"`
	Petid           *int   `json:"petid,omitempty"`
	Flag            *int   `json:"flag,omitempty"`
	Expiration      int64  `json:"expiration,omitempty"`
	GiftFrom        string `gorm:"column:giftFrom" json:"giftFrom,omitempty"`
}

func (Inventoryitem) TableName() string {
	return "inventoryitems"
}
