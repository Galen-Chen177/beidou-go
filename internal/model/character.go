package model

import "time"

// Character 角色模型
type Character struct {
	ID       int32  `gorm:"primaryKey;column:id"`
	AccountID int32 `gorm:"column:accountid;index"`
	WorldID  int8  `gorm:"column:world"`
	Name     string `gorm:"column:name;uniqueIndex;size:13"`

	Level    int8   `gorm:"column:level"`
	Job      int16  `gorm:"column:job"`
	Str      int16  `gorm:"column:str"`
	Dex      int16  `gorm:"column:dex"`
	Int      int16  `gorm:"column:int"`
	Luk      int16  `gorm:"column:luk"`
	HP       int16  `gorm:"column:hp"`
	MaxHP    int16  `gorm:"column:maxhp"`
	MP       int16  `gorm:"column:mp"`
	MaxMP    int16  `gorm:"column:maxmp"`
	AP       int16  `gorm:"column:ap"`
	SP       int16  `gorm:"column:sp"`
	EXP      int32  `gorm:"column:exp"`
	Fame     int16  `gorm:"column:fame"`
	Meso     int32  `gorm:"column:meso"`

	MapID     int32 `gorm:"column:map"`
	SpawnPoint int8 `gorm:"column:spawnpoint"`

	Gender byte   `gorm:"column:gender"`
	Skin   byte   `gorm:"column:skin"`
	Face   int32  `gorm:"column:face"`
	Hair   int32  `gorm:"column:hair"`

	CreatedAt time.Time `gorm:"column:createdat"`
	UpdatedAt time.Time `gorm:"column:updatedat"`
}

func (Character) TableName() string {
	return "characters"
}
