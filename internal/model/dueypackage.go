package model

import (
	"time"

	"gorm.io/gorm"
)

// Dueypackage 对应表 dueypackages
type Dueypackage struct {
	gorm.Model

	Packageid  int64     `gorm:"primaryKey;column:packageid" json:"packageid,omitempty"`
	Receiverid int64     `gorm:"column:receiverid" json:"receiverid,omitempty"`
	Sendername string    `gorm:"column:sendername;type:varchar(200)" json:"sendername,omitempty"`
	Mesos      int64     `gorm:"column:mesos" json:"mesos,omitempty"`
	Timestamp  time.Time `gorm:"column:timestamp" json:"timestamp,omitempty"`
	Message    string    `gorm:"column:message;type:varchar(200)" json:"message,omitempty"`
	Checked    int       `gorm:"column:checked" json:"checked,omitempty"`
	Type       int       `gorm:"column:type" json:"type,omitempty"`
}

func (Dueypackage) TableName() string {
	return "dueypackages"
}
