package model

// FamilyEntitlement 实体映射
type FamilyEntitlement struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	Charid        *int   `json:"charid,omitempty"`
	Entitlementid *int   `json:"entitlementid,omitempty"`
	Timestamp     int64  `json:"timestamp,omitempty"`
}

func (FamilyEntitlement) TableName() string {
	return "family_entitlement"
}
