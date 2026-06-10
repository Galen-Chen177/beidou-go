package model

import (
	"time"

	"gorm.io/gorm"
)

// Allianceguild 对应表 allianceguilds
type Allianceguild struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Allianceid int `gorm:"column:allianceid" json:"allianceid,omitempty"`
	Guildid    int `gorm:"column:guildid" json:"guildid,omitempty"`
}

func (Allianceguild) TableName() string {
	return "allianceguilds"
}
