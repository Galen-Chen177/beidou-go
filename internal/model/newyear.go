package model

import (
	"gorm.io/gorm"
)

// Newyear 映射 newyear 表
type Newyear struct {
	gorm.Model
	Senderid        int    `json:"senderid,omitempty"`
	Sendername      string `gorm:"column:sendername;type:varchar(200)" json:"sendername,omitempty"`
	Receiverid      int    `json:"receiverid,omitempty"`
	Receivername    string `gorm:"column:receivername;type:varchar(200)" json:"receivername,omitempty"`
	Message         string `gorm:"column:message;type:varchar(200)" json:"message,omitempty"`
	Senderdiscard   bool   `json:"senderdiscard,omitempty"`
	Receiverdiscard bool   `json:"receiverdiscard,omitempty"`
	Received        bool   `json:"received,omitempty"`
	Timesent        int64  `json:"timesent,omitempty"`
	Timereceived    int64  `json:"timereceived,omitempty"`
}

// TableName 指定表名
func (Newyear) TableName() string {
	return "newyear"
}
