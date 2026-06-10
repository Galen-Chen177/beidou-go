package model

import (
	"time"

	"gorm.io/gorm"
)

// Makerrewarddata 对应表 makerrewarddata (复合主键：itemid + rewardid)
type Makerrewarddata struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Itemid    int            `gorm:"column:itemid" json:"itemid"`
	Rewardid  int            `gorm:"column:rewardid" json:"rewardid"`
	Quantity  int            `gorm:"column:quantity" json:"quantity"`
	Prob      int            `gorm:"column:prob" json:"prob"`
}

func (Makerrewarddata) TableName() string {
	return "makerrewarddata"
}
