package model

import (
	"time"

	"gorm.io/gorm"
)

// Namechange 映射 namechanges 表
type Namechange struct {
	gorm.Model
	Characterid    int       `json:"characterid,omitempty"`
	Older          string    `gorm:"column:old;type:varchar(200)" json:"older,omitempty"`
	Newer          string    `gorm:"column:new;type:varchar(200)" json:"newer,omitempty"`
	RequestTime    time.Time `gorm:"column:requestTime" json:"requestTime,omitempty"`
	CompletionTime time.Time `gorm:"column:completionTime" json:"completionTime,omitempty"`
}

// TableName 指定表名
func (Namechange) TableName() string {
	return "namechanges"
}
