package model

import (
	"time"

	"gorm.io/gorm"
)

// LangResource 对应表 lang_resources
type LangResource struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	LangType   string         `gorm:"column:langType" json:"langType"`
	LangBase   string         `gorm:"column:langBase" json:"langBase"`
	LangCode   string         `gorm:"column:langCode" json:"langCode"`
	LangValue  string         `gorm:"column:langValue" json:"langValue"`
	LangExtend string         `gorm:"column:langExtend" json:"langExtend"`
}

func (LangResource) TableName() string {
	return "lang_resources"
}
