package model

import (
	"gorm.io/gorm"
)

// Response mapped from table "responses"
type Response struct {
	gorm.Model
	Chat     string `gorm:"column:chat;type:varchar(200)" json:"chat,omitempty"`
	Response string `gorm:"column:response;type:varchar(200)" json:"response,omitempty"`
}

func (Response) TableName() string {
	return "responses"
}
