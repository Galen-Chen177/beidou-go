package model

import (
	"time"

	"gorm.io/gorm"
)

// Account 对应表 accounts
type Account struct {
	gorm.Model

	Name           string    `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	Password       string    `gorm:"column:password;type:varchar(200)" json:"-"`
	Pin            string    `gorm:"column:pin;type:varchar(200)" json:"pin,omitempty"`
	Pic            string    `gorm:"column:pic;type:varchar(200)" json:"pic,omitempty"`
	Loggedin       int       `gorm:"column:loggedin" json:"loggedin,omitempty"`
	Lastlogin      time.Time `gorm:"column:lastlogin" json:"lastlogin,omitempty"`
	Birthday       time.Time `gorm:"column:birthday" json:"birthday,omitempty"`
	Banned         bool      `gorm:"column:banned" json:"banned,omitempty"`
	Banreason      string    `gorm:"column:banreason;type:varchar(200)" json:"banreason,omitempty"`
	Macs           string    `gorm:"column:macs;type:varchar(200)" json:"macs,omitempty"`
	NxCredit       int       `gorm:"column:nxCredit" json:"nxCredit,omitempty"`
	MaplePoint     int       `gorm:"column:maplePoint" json:"maplePoint,omitempty"`
	NxPrepaid      int       `gorm:"column:nxPrepaid" json:"nxPrepaid,omitempty"`
	Characterslots int       `gorm:"column:characterslots" json:"characterslots,omitempty"`
	Gender         int       `gorm:"column:gender" json:"gender,omitempty"`
	Tempban        time.Time `gorm:"column:tempban" json:"tempban,omitempty"`
	Greason        int       `gorm:"column:greason" json:"greason,omitempty"`
	Tos            bool      `gorm:"column:tos" json:"tos,omitempty"`
	Sitelogged     string    `gorm:"column:sitelogged;type:varchar(200)" json:"sitelogged,omitempty"`
	Webadmin       int       `gorm:"column:webadmin" json:"webadmin,omitempty"`
	Nick           string    `gorm:"column:nick;type:varchar(200)" json:"nick,omitempty"`
	Mute           int       `gorm:"column:mute" json:"mute,omitempty"`
	Email          string    `gorm:"column:email;type:varchar(200)" json:"email,omitempty"`
	IP             string    `gorm:"column:ip;type:varchar(200)" json:"ip,omitempty"`
	Rewardpoints   int       `gorm:"column:rewardpoints" json:"rewardpoints,omitempty"`
	Votepoints     int       `gorm:"column:votepoints" json:"votepoints,omitempty"`
	Hwid           string    `gorm:"column:hwid;type:varchar(200)" json:"hwid,omitempty"`
	Language       int       `gorm:"column:language" json:"language,omitempty"`
}

func (Account) TableName() string {
	return "accounts"
}
