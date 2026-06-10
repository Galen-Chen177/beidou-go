package model

// HpMpAlert 实体映射
type HpMpAlert struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	CId *int  `gorm:"column:cId" json:"cId,omitempty"`
	Hp  int8  `json:"hp,omitempty"`
	Mp  int8  `json:"mp,omitempty"`
}

func (HpMpAlert) TableName() string {
	return "hp_mp_alert"
}
