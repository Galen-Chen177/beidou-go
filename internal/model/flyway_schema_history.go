package model

import "time"

// FlywaySchemaHistory 实体映射
type FlywaySchemaHistory struct {
	InstalledRank int       `gorm:"primaryKey;column:installed_rank" json:"installedRank,omitempty"`
	Version       string    `json:"version,omitempty"`
	Description   string    `json:"description,omitempty"`
	Type          string    `json:"type,omitempty"`
	Script        string    `json:"script,omitempty"`
	Checksum      *int      `json:"checksum,omitempty"`
	InstalledBy   string    `json:"installedBy,omitempty"`
	InstalledOn   time.Time `json:"installedOn,omitempty"`
	ExecutionTime *int      `json:"executionTime,omitempty"`
	Success       *bool     `json:"success,omitempty"`
}

func (FlywaySchemaHistory) TableName() string {
	return "flyway_schema_history"
}
