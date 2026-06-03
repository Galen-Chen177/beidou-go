package model

import "time"

// Account 账号模型
type Account struct {
	ID        int32     `gorm:"primaryKey;column:id"`
	Name      string    `gorm:"column:name;uniqueIndex;size:13"`
	Password  string    `gorm:"column:password;size:128"` // SHA-512 hash
	Pin       string    `gorm:"column:pin;size:10"`
	Pic       string    `gorm:"column:pic;size:10"`
	GMLevel   int8      `gorm:"column:gm"`
	Banned    bool      `gorm:"column:banned"`
	CreatedAt time.Time `gorm:"column:createdat"`
	LastLogin time.Time `gorm:"column:lastlogin"`
}

func (Account) TableName() string {
	return "accounts"
}
