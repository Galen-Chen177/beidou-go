package model

// Hwidban 实体映射
type Hwidban struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	Hwidbanid int64  `gorm:"column:hwidbanid" json:"hwidbanid,omitempty"`
	Hwid      string `json:"hwid,omitempty"`
}

func (Hwidban) TableName() string {
	return "hwidbans"
}
