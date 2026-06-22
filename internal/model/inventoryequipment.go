package model

import "gorm.io/gorm"

// Inventoryequipment 实体映射
type Inventoryequipment struct {
	gorm.Model
	Inventoryequipmentid int64 `gorm:"column:inventoryequipmentid" json:"inventoryequipmentid,omitempty"`
	Inventoryitemid      int64 `json:"inventoryitemid,omitempty"`
	Upgradeslots         *int  `json:"upgradeslots,omitempty"`
	Level                *int  `json:"level,omitempty"`
	Str                  *int  `json:"str,omitempty"`
	Dex                  *int  `json:"dex,omitempty"`
	Inte                 *int  `gorm:"column:int" json:"inte,omitempty"`
	Luk                  *int  `json:"luk,omitempty"`
	Hp                   *int  `json:"hp,omitempty"`
	Mp                   *int  `json:"mp,omitempty"`
	Watk                 *int  `json:"watk,omitempty"`
	Matk                 *int  `json:"matk,omitempty"`
	Wdef                 *int  `json:"wdef,omitempty"`
	Mdef                 *int  `json:"mdef,omitempty"`
	Acc                  *int  `json:"acc,omitempty"`
	Avoid                *int  `json:"avoid,omitempty"`
	Hands                *int  `json:"hands,omitempty"`
	Speed                *int  `json:"speed,omitempty"`
	Jump                 *int  `json:"jump,omitempty"`
	Locked               *int  `json:"locked,omitempty"`
	Vicious              *int  `json:"vicious,omitempty"`
	Itemlevel            *int  `json:"itemlevel,omitempty"`
	Itemexp              *int  `json:"itemexp,omitempty"`
	Ringid               *int  `json:"ringid,omitempty"`
}

func (Inventoryequipment) TableName() string {
	return "inventoryequipment"
}
