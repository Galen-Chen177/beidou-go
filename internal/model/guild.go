package model

// Guild 实体映射
type Guild struct {
	Guildid     int64  `gorm:"primaryKey;column:guildid" json:"guildid,omitempty"`
	Leader      int64  `json:"leader,omitempty"`
	Gp          int64  `json:"gp,omitempty"`
	Logo        int64  `json:"logo,omitempty"`
	LogoColor   *int   `gorm:"column:logoColor" json:"logoColor,omitempty"`
	Name        string `gorm:"column:name;type:varchar(200)" json:"name,omitempty"`
	Rank1title  string `gorm:"column:rank1title;type:varchar(200)" json:"rank1title,omitempty"`
	Rank2title  string `gorm:"column:rank2title;type:varchar(200)" json:"rank2title,omitempty"`
	Rank3title  string `gorm:"column:rank3title;type:varchar(200)" json:"rank3title,omitempty"`
	Rank4title  string `gorm:"column:rank4title;type:varchar(200)" json:"rank4title,omitempty"`
	Rank5title  string `gorm:"column:rank5title;type:varchar(200)" json:"rank5title,omitempty"`
	Capacity    int64  `json:"capacity,omitempty"`
	LogoBG      int64  `gorm:"column:logoBG" json:"logoBG,omitempty"`
	LogoBGColor *int   `gorm:"column:logoBGColor" json:"logoBGColor,omitempty"`
	Notice      string `gorm:"column:notice;type:varchar(200)" json:"notice,omitempty"`
	Signature   *int   `json:"signature,omitempty"`
	AllianceId  int64  `gorm:"column:allianceId" json:"allianceId,omitempty"`
}

func (Guild) TableName() string {
	return "guilds"
}
