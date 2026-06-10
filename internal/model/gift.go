package model

// Gift 实体映射
type Gift struct {
	ID        uint64         `gorm:"primaryKey;type:bigint;comment:自增长id"`
	To      *int   `json:"to,omitempty"`
	From    string  `json:"from,omitempty"`
	Message string  `json:"message,omitempty"`
	Sn      int64   `json:"sn,omitempty"`
	Ringid  *int    `json:"ringid,omitempty"`
}

func (Gift) TableName() string {
	return "gifts"
}
