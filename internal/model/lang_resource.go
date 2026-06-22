package model

import (
	"gorm.io/gorm"
)

// LangResource 对应表 lang_resources
type LangResource struct {
	gorm.Model
	LangType   string `gorm:"column:langType;type:varchar(200)" json:"langType"`
	LangBase   string `gorm:"column:langBase;type:varchar(200)" json:"langBase"`
	LangCode   string `gorm:"column:langCode;type:varchar(200)" json:"langCode"`
	LangValue  string `gorm:"column:langValue;type:varchar(200)" json:"langValue"`
	LangExtend string `gorm:"column:langExtend;type:varchar(200)" json:"langExtend"`
}

func (LangResource) TableName() string {
	return "lang_resources"
}
