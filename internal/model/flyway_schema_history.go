package model

import "time"

// FlywaySchemaHistory 实体映射
type FlywaySchemaHistory struct {
	InstalledRank int       `gorm:"primaryKey;column:installed_rank" json:"installedRank,omitempty"`
	Version       string    `gorm:"column:version;type:varchar(200)" json:"version,omitempty"`
	Description   string    `gorm:"column:description;type:varchar(200)" json:"description,omitempty"`
	Type          string    `gorm:"column:type;type:varchar(200)" json:"type,omitempty"`
	Script        string    `gorm:"column:script;type:varchar(200)" json:"script,omitempty"`
	Checksum      *int      `json:"checksum,omitempty"`
	InstalledBy   string    `gorm:"column:installedBy;type:varchar(200)" json:"installedBy,omitempty"`
	InstalledOn   time.Time `json:"installedOn,omitempty"`
	ExecutionTime *int      `json:"executionTime,omitempty"`
	Success       *bool     `json:"success,omitempty"`
}

func (FlywaySchemaHistory) TableName() string {
	return "flyway_schema_history"
}
