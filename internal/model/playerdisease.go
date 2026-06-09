package model

import (
	"time"

	"gorm.io/gorm"
)

// Playerdisease 映射 playerdiseases 表
type Playerdisease struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Charid     int            `json:"charid,omitempty"`
	Disease    int            `json:"disease,omitempty"`
	Mobskillid int            `json:"mobskillid,omitempty"`
	Mobskilllv int            `json:"mobskilllv,omitempty"`
	Length     int64          `json:"length,omitempty"`
}

// TableName 指定表名
func (Playerdisease) TableName() string {
	return "playerdiseases"
}
