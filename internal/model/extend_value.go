package model

// ExtendValue 扩展字段表 实体映射
type ExtendValue struct {
	ExtendId    string `gorm:"primaryKey;column:extend_id" json:"extendId,omitempty"`
	ExtendType  string `gorm:"primaryKey;column:extend_type" json:"extendType,omitempty"`
	ExtendName  string `gorm:"primaryKey;column:extend_name" json:"extendName,omitempty"`
	ExtendValue string `json:"extendValue,omitempty"`
	CreateTime  string `json:"createTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
}

func (ExtendValue) TableName() string {
	return "extend_value"
}
