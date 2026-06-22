package model

import "gorm.io/gorm"

// FamilyEntitlement 实体映射
type FamilyEntitlement struct {
	gorm.Model
	Charid        *int  `json:"charid,omitempty"`
	Entitlementid *int  `json:"entitlementid,omitempty"`
	Timestamp     int64 `json:"timestamp,omitempty"`
}

func (FamilyEntitlement) TableName() string {
	return "family_entitlement"
}
