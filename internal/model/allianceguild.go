package model

import (
	"gorm.io/gorm"
)

// Allianceguild 对应表 allianceguilds
type Allianceguild struct {
	gorm.Model

	Allianceid int `gorm:"column:allianceid" json:"allianceid,omitempty"`
	Guildid    int `gorm:"column:guildid" json:"guildid,omitempty"`
}

func (Allianceguild) TableName() string {
	return "allianceguilds"
}
