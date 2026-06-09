package model

import (
	"time"

	"gorm.io/gorm"
)

// Worldtransfer mapped from table "worldtransfers"
type Worldtransfer struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	TransferID     int            `gorm:"column:id;autoIncrement" json:"transferID,omitempty"`
	Characterid    int            `gorm:"column:characterid" json:"characterid,omitempty"`
	From           int            `gorm:"column:from" json:"from,omitempty"`
	To             int            `gorm:"column:to" json:"to,omitempty"`
	RequestTime    time.Time      `gorm:"column:requestTime" json:"requestTime,omitempty"`
	CompletionTime time.Time      `gorm:"column:completionTime" json:"completionTime,omitempty"`
}

func (Worldtransfer) TableName() string {
	return "worldtransfers"
}
