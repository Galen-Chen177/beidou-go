package model

// Eventstat 实体映射
type Eventstat struct {
	Characterid int64  `gorm:"primaryKey" json:"characterid,omitempty"`
	Name        string `json:"name,omitempty"`
	Info        *int   `json:"info,omitempty"`
}

func (Eventstat) TableName() string {
	return "eventstats"
}
