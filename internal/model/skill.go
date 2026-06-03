package model

// Skill 角色技能模型
type Skill struct {
	ID          int64 `gorm:"primaryKey;column:id;autoIncrement"`
	CharacterID int32 `gorm:"column:characterid;index"`
	SkillID     int32 `gorm:"column:skillid"`
	Level       int8  `gorm:"column:level"`
	MasterLevel int8  `gorm:"column:masterlevel"`
}

func (Skill) TableName() string {
	return "skills"
}
