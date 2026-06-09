package model

import (
	"math/big"
	"time"

	"gorm.io/gorm"
)

// Nxcode 映射 nxcode 表
type Nxcode struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Code       string         `json:"code,omitempty"`
	Retriever  string         `json:"retriever,omitempty"`
	Expiration *big.Int       `json:"expiration,omitempty"`
}

// TableName 指定表名
func (Nxcode) TableName() string {
	return "nxcode"
}
