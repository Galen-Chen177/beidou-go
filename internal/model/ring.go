package model

import (
	"time"

	"gorm.io/gorm"
)

// Ring mapped from table "rings"
type Ring struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	RingID        int            `gorm:"column:id;autoIncrement" json:"ringID,omitempty"`
	PartnerRingId int            `gorm:"column:partnerRingId" json:"partnerRingId,omitempty"`
	PartnerChrId  int            `gorm:"column:partnerChrId" json:"partnerChrId,omitempty"`
	Itemid        int            `gorm:"column:itemid" json:"itemid,omitempty"`
	Partnername   string         `gorm:"column:partnername" json:"partnername,omitempty"`
}

func (Ring) TableName() string {
	return "rings"
}
