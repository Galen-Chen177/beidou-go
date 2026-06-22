package model

import (
	"gorm.io/gorm"
)

// Ring mapped from table "rings"
type Ring struct {
	gorm.Model
	PartnerRingId int    `gorm:"column:partnerRingId" json:"partnerRingId,omitempty"`
	PartnerChrId  int    `gorm:"column:partnerChrId" json:"partnerChrId,omitempty"`
	Itemid        int    `gorm:"column:itemid" json:"itemid,omitempty"`
	Partnername   string `gorm:"column:partnername;type:varchar(200)" json:"partnername,omitempty"`
}

func (Ring) TableName() string {
	return "rings"
}
