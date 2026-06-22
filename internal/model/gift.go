package model

import "gorm.io/gorm"

// Gift 实体映射
type Gift struct {
	gorm.Model
	To      *int   `json:"to,omitempty"`
	From    string `gorm:"column:from;type:varchar(200)" json:"from,omitempty"`
	Message string `gorm:"column:message;type:varchar(200)" json:"message,omitempty"`
	Sn      int64  `json:"sn,omitempty"`
	Ringid  *int   `json:"ringid,omitempty"`
}

func (Gift) TableName() string {
	return "gifts"
}
