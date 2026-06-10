package model

import (
	"time"

	"gorm.io/gorm"
)

// Account 对应表 accounts
type Account struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CreatedAt *time.Time     `json:"createdAt"`
	UpdatedAt *time.Time     `json:"-"`
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Name           string    `gorm:"column:name" json:"name,omitempty"`
	Password       string    `gorm:"column:password" json:"-"`
	Pin            string    `gorm:"column:pin" json:"pin,omitempty"`
	Pic            string    `gorm:"column:pic" json:"pic,omitempty"`
	Loggedin       int       `gorm:"column:loggedin" json:"loggedin,omitempty"`
	Lastlogin      time.Time `gorm:"column:lastlogin" json:"lastlogin,omitempty"`
	Createdat      time.Time `gorm:"column:createdat" json:"createdat,omitempty"`
	Birthday       time.Time `gorm:"column:birthday" json:"birthday,omitempty"`
	Banned         bool      `gorm:"column:banned" json:"banned,omitempty"`
	Banreason      string    `gorm:"column:banreason" json:"banreason,omitempty"`
	Macs           string    `gorm:"column:macs" json:"macs,omitempty"`
	NxCredit       int       `gorm:"column:nxCredit" json:"nxCredit,omitempty"`
	MaplePoint     int       `gorm:"column:maplePoint" json:"maplePoint,omitempty"`
	NxPrepaid      int       `gorm:"column:nxPrepaid" json:"nxPrepaid,omitempty"`
	Characterslots int       `gorm:"column:characterslots" json:"characterslots,omitempty"`
	Gender         int       `gorm:"column:gender" json:"gender,omitempty"`
	Tempban        time.Time `gorm:"column:tempban" json:"tempban,omitempty"`
	Greason        int       `gorm:"column:greason" json:"greason,omitempty"`
	Tos            bool      `gorm:"column:tos" json:"tos,omitempty"`
	Sitelogged     string    `gorm:"column:sitelogged" json:"sitelogged,omitempty"`
	Webadmin       int       `gorm:"column:webadmin" json:"webadmin,omitempty"`
	Nick           string    `gorm:"column:nick" json:"nick,omitempty"`
	Mute           int       `gorm:"column:mute" json:"mute,omitempty"`
	Email          string    `gorm:"column:email" json:"email,omitempty"`
	IP             string    `gorm:"column:ip" json:"ip,omitempty"`
	Rewardpoints   int       `gorm:"column:rewardpoints" json:"rewardpoints,omitempty"`
	Votepoints     int       `gorm:"column:votepoints" json:"votepoints,omitempty"`
	Hwid           string    `gorm:"column:hwid" json:"hwid,omitempty"`
	Language       int       `gorm:"column:language" json:"language,omitempty"`
}

func (Account) TableName() string {
	return "accounts"
}
