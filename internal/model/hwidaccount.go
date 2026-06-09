package model

import "time"

// Hwidaccount 实体映射
type Hwidaccount struct {
	Accountid int       `gorm:"primaryKey;column:accountid" json:"accountid,omitempty"`
	Hwid      string    `gorm:"primaryKey;column:hwid" json:"hwid,omitempty"`
	Relevance *int      `json:"relevance,omitempty"`
	Expiresat time.Time `json:"expiresat,omitempty"`
}

func (Hwidaccount) TableName() string {
	return "hwidaccounts"
}
