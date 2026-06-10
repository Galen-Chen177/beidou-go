package model

import (
	"time"

	"gorm.io/gorm"
)

// Response mapped from table "responses"
type Response struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time `json:"createdAt"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Chat      string         `gorm:"column:chat" json:"chat,omitempty"`
	Response  string         `gorm:"column:response" json:"response,omitempty"`
	ResID     int64          `gorm:"column:id;autoIncrement" json:"resID,omitempty"`
}

func (Response) TableName() string {
	return "responses"
}
