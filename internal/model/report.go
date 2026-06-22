package model

import (
	"time"

	"gorm.io/gorm"
)

// Report mapped from table "reports"
type Report struct {
	gorm.Model
	Reporttime  time.Time `gorm:"column:reporttime" json:"reporttime,omitempty"`
	Reporterid  int       `gorm:"column:reporterid" json:"reporterid,omitempty"`
	Victimid    int       `gorm:"column:victimid" json:"victimid,omitempty"`
	Reason      int       `gorm:"column:reason" json:"reason,omitempty"`
	Chatlog     string    `gorm:"column:chatlog;type:varchar(200)" json:"chatlog,omitempty"`
	Description string    `gorm:"column:description;type:varchar(200)" json:"description,omitempty"`
}

func (Report) TableName() string {
	return "reports"
}
