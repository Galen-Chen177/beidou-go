package model

import (
	"time"

	"gorm.io/gorm"
)

// Dueypackage 对应表 dueypackages
type Dueypackage struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Packageid  int64     `gorm:"primaryKey;column:packageid" json:"packageid,omitempty"`
	Receiverid int64     `gorm:"column:receiverid" json:"receiverid,omitempty"`
	Sendername string    `gorm:"column:sendername" json:"sendername,omitempty"`
	Mesos      int64     `gorm:"column:mesos" json:"mesos,omitempty"`
	Timestamp  time.Time `gorm:"column:timestamp" json:"timestamp,omitempty"`
	Message    string    `gorm:"column:message" json:"message,omitempty"`
	Checked    int       `gorm:"column:checked" json:"checked,omitempty"`
	Type       int       `gorm:"column:type" json:"type,omitempty"`
}

func (Dueypackage) TableName() string {
	return "dueypackages"
}
