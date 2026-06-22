package model

import (
	"time"

	"gorm.io/gorm"
)

// Worldtransfer mapped from table "worldtransfers"
type Worldtransfer struct {
	gorm.Model
	Characterid    int       `gorm:"column:characterid" json:"characterid,omitempty"`
	From           int       `gorm:"column:from" json:"from,omitempty"`
	To             int       `gorm:"column:to" json:"to,omitempty"`
	RequestTime    time.Time `gorm:"column:requestTime" json:"requestTime,omitempty"`
	CompletionTime time.Time `gorm:"column:completionTime" json:"completionTime,omitempty"`
}

func (Worldtransfer) TableName() string {
	return "worldtransfers"
}
