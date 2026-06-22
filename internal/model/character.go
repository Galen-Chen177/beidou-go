package model

import (
	"time"

	"gorm.io/gorm"
)

// Character 对应表 characters
type Character struct {
	gorm.Model

	Accountid            int       `gorm:"column:accountid" json:"accountid,omitempty"`
	World                int       `gorm:"column:world" json:"world,omitempty"`
	Name                 string    `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	Level                int       `gorm:"column:level" json:"level,omitempty"`
	Exp                  int       `gorm:"column:exp" json:"exp,omitempty"`
	Gachaexp             int       `gorm:"column:gachaexp" json:"gachaexp,omitempty"`
	AttrStr              int       `gorm:"column:str" json:"attrStr,omitempty"`
	AttrDex              int       `gorm:"column:dex" json:"attrDex,omitempty"`
	AttrLuk              int       `gorm:"column:luk" json:"attrLuk,omitempty"`
	AttrInt              int       `gorm:"column:int" json:"attrInt,omitempty"`
	Hp                   int       `gorm:"column:hp" json:"hp,omitempty"`
	Mp                   int       `gorm:"column:mp" json:"mp,omitempty"`
	Maxhp                int       `gorm:"column:maxhp" json:"maxhp,omitempty"`
	Maxmp                int       `gorm:"column:maxmp" json:"maxmp,omitempty"`
	Meso                 int       `gorm:"column:meso" json:"meso,omitempty"`
	HpMpUsed             int       `gorm:"column:hpMpUsed" json:"hpMpUsed,omitempty"`
	Job                  int       `gorm:"column:job" json:"job,omitempty"`
	Skincolor            int       `gorm:"column:skincolor" json:"skincolor,omitempty"`
	Gender               int       `gorm:"column:gender" json:"gender,omitempty"`
	Fame                 int       `gorm:"column:fame" json:"fame,omitempty"`
	Fquest               int       `gorm:"column:fquest" json:"fquest,omitempty"`
	Hair                 int       `gorm:"column:hair" json:"hair,omitempty"`
	Face                 int       `gorm:"column:face" json:"face,omitempty"`
	Ap                   int       `gorm:"column:ap" json:"ap,omitempty"`
	Sp                   string    `gorm:"column:sp;type:varchar(200)" json:"sp,omitempty"`
	Map                  int       `gorm:"column:map" json:"map,omitempty"`
	Spawnpoint           int       `gorm:"column:spawnpoint" json:"spawnpoint,omitempty"`
	Gm                   int       `gorm:"column:gm" json:"gm,omitempty"`
	Party                int       `gorm:"column:party" json:"party,omitempty"`
	BuddyCapacity        int       `gorm:"column:buddyCapacity" json:"buddyCapacity,omitempty"`
	Createdate           time.Time `gorm:"column:createdate" json:"createdate,omitempty"`
	Rank                 int       `gorm:"column:rank" json:"rank,omitempty"`
	RankMove             int       `gorm:"column:rankMove" json:"rankMove,omitempty"`
	JobRank              int       `gorm:"column:jobRank" json:"jobRank,omitempty"`
	JobRankMove          int       `gorm:"column:jobRankMove" json:"jobRankMove,omitempty"`
	Guildid              int       `gorm:"column:guildid" json:"guildid,omitempty"`
	Guildrank            int       `gorm:"column:guildrank" json:"guildrank,omitempty"`
	Messengerid          int       `gorm:"column:messengerid" json:"messengerid,omitempty"`
	Messengerposition    int       `gorm:"column:messengerposition" json:"messengerposition,omitempty"`
	Mountlevel           int       `gorm:"column:mountlevel" json:"mountlevel,omitempty"`
	Mountexp             int       `gorm:"column:mountexp" json:"mountexp,omitempty"`
	Mounttiredness       int       `gorm:"column:mounttiredness" json:"mounttiredness,omitempty"`
	Omokwins             int       `gorm:"column:omokwins" json:"omokwins,omitempty"`
	Omoklosses           int       `gorm:"column:omoklosses" json:"omoklosses,omitempty"`
	Omokties             int       `gorm:"column:omokties" json:"omokties,omitempty"`
	Matchcardwins        int       `gorm:"column:matchcardwins" json:"matchcardwins,omitempty"`
	Matchcardlosses      int       `gorm:"column:matchcardlosses" json:"matchcardlosses,omitempty"`
	Matchcardties        int       `gorm:"column:matchcardties" json:"matchcardties,omitempty"`
	Merchantmesos        int       `gorm:"column:merchantmesos" json:"merchantmesos,omitempty"`
	Hasmerchant          bool      `gorm:"column:hasmerchant" json:"hasmerchant,omitempty"`
	Equipslots           int       `gorm:"column:equipslots" json:"equipslots,omitempty"`
	Useslots             int       `gorm:"column:useslots" json:"useslots,omitempty"`
	Setupslots           int       `gorm:"column:setupslots" json:"setupslots,omitempty"`
	Etcslots             int       `gorm:"column:etcslots" json:"etcslots,omitempty"`
	FamilyId             int       `gorm:"column:familyId" json:"familyId,omitempty"`
	Monsterbookcover     int       `gorm:"column:monsterbookcover" json:"monsterbookcover,omitempty"`
	AllianceRank         int       `gorm:"column:allianceRank" json:"allianceRank,omitempty"`
	VanquisherStage      int       `gorm:"column:vanquisherStage" json:"vanquisherStage,omitempty"`
	AriantPoints         int       `gorm:"column:ariantPoints" json:"ariantPoints,omitempty"`
	DojoPoints           int       `gorm:"column:dojoPoints" json:"dojoPoints,omitempty"`
	LastDojoStage        int       `gorm:"column:lastDojoStage" json:"lastDojoStage,omitempty"`
	FinishedDojoTutorial int       `gorm:"column:finishedDojoTutorial" json:"finishedDojoTutorial,omitempty"`
	VanquisherKills      int       `gorm:"column:vanquisherKills" json:"vanquisherKills,omitempty"`
	SummonValue          int64     `gorm:"column:summonValue" json:"summonValue,omitempty"`
	PartnerId            int       `gorm:"column:partnerId" json:"partnerId,omitempty"`
	MarriageItemId       int       `gorm:"column:marriageItemId" json:"marriageItemId,omitempty"`
	Reborns              int       `gorm:"column:reborns" json:"reborns,omitempty"`
	Pqpoints             int       `gorm:"column:pqpoints" json:"pqpoints,omitempty"`
	DataString           string    `gorm:"column:dataString;type:varchar(200)" json:"dataString,omitempty"`
	LastLogoutTime       time.Time `gorm:"column:lastLogoutTime" json:"lastLogoutTime,omitempty"`
	LastExpGainTime      time.Time `gorm:"column:lastExpGainTime" json:"lastExpGainTime,omitempty"`
	PartySearch          bool      `gorm:"column:partySearch" json:"partySearch,omitempty"`
	Jailexpire           int64     `gorm:"column:jailexpire" json:"jailexpire,omitempty"`
}

func (Character) TableName() string {
	return "characters"
}
