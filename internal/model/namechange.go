package model

import (
	"time"

	"gorm.io/gorm"
)

// Namechange 映射 namechanges 表
type Namechange struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Characterid    int            `json:"characterid,omitempty"`
	Older          string         `gorm:"column:old" json:"older,omitempty"`
	Newer          string         `gorm:"column:new" json:"newer,omitempty"`
	RequestTime    time.Time      `gorm:"column:requestTime" json:"requestTime,omitempty"`
	CompletionTime time.Time      `gorm:"column:completionTime" json:"completionTime,omitempty"`
}

// TableName 指定表名
func (Namechange) TableName() string {
	return "namechanges"
}
