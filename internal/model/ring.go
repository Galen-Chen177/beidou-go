package model

import (
	"gorm.io/gorm"
)

// Ring mapped from table "rings"
type Ring struct {
	gorm.Model
	PartnerRingID int    `gorm:"column:partnerRingID" json:"partnerRingID,omitempty"`
	PartnerChrID  int    `gorm:"column:partnerChrID" json:"partnerChrID,omitempty"`
	Itemid        int    `gorm:"column:itemid" json:"itemid,omitempty"`
	Partnername   string `gorm:"column:partnername;type:varchar(200)" json:"partnername,omitempty"`
}

func (Ring) TableName() string {
	return "rings"
}
