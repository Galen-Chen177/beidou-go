package model

import (
	"time"

	"gorm.io/gorm"
)

// Report mapped from table "reports"
type Report struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Reportid    int64          `gorm:"column:id;autoIncrement" json:"reportid,omitempty"`
	Reporttime  time.Time      `gorm:"column:reporttime" json:"reporttime,omitempty"`
	Reporterid  int            `gorm:"column:reporterid" json:"reporterid,omitempty"`
	Victimid    int            `gorm:"column:victimid" json:"victimid,omitempty"`
	Reason      int            `gorm:"column:reason" json:"reason,omitempty"`
	Chatlog     string         `gorm:"column:chatlog" json:"chatlog,omitempty"`
	Description string         `gorm:"column:description" json:"description,omitempty"`
}

func (Report) TableName() string {
	return "reports"
}
