package model

// Guild 实体映射
type Guild struct {
	Guildid    int64   `gorm:"primaryKey;column:guildid" json:"guildid,omitempty"`
	Leader     int64    `json:"leader,omitempty"`
	Gp         int64    `json:"gp,omitempty"`
	Logo       int64    `json:"logo,omitempty"`
	LogoColor  *int     `gorm:"column:logoColor" json:"logoColor,omitempty"`
	Name       string   `json:"name,omitempty"`
	Rank1title string   `json:"rank1title,omitempty"`
	Rank2title string   `json:"rank2title,omitempty"`
	Rank3title string   `json:"rank3title,omitempty"`
	Rank4title string   `json:"rank4title,omitempty"`
	Rank5title string   `json:"rank5title,omitempty"`
	Capacity   int64    `json:"capacity,omitempty"`
	LogoBG     int64    `gorm:"column:logoBG" json:"logoBG,omitempty"`
	LogoBGColor *int    `gorm:"column:logoBGColor" json:"logoBGColor,omitempty"`
	Notice     string   `json:"notice,omitempty"`
	Signature  *int     `json:"signature,omitempty"`
	AllianceId int64    `gorm:"column:allianceId" json:"allianceId,omitempty"`
}

func (Guild) TableName() string {
	return "guilds"
}
