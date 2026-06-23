package model

import "gorm.io/gorm"

// ExtendValue 扩展字段表 实体映射
type ExtendValue struct {
	gorm.Model
	ExtendID    string `gorm:"primaryKey;column:extend_id;type:varchar(200)" json:"extendID,omitempty"`
	ExtendType  string `gorm:"primaryKey;column:extend_type;type:varchar(200)" json:"extendType,omitempty"`
	ExtendName  string `gorm:"primaryKey;column:extend_name;type:varchar(200)" json:"extendName,omitempty"`
	ExtendValue string `gorm:"column:extendValue;type:varchar(200)" json:"extendValue,omitempty"`
	CreateTime  string `gorm:"column:createTime;type:varchar(200)" json:"createTime,omitempty"`
	UpdateTime  string `gorm:"column:updateTime;type:varchar(200)" json:"updateTime,omitempty"`
}

func (ExtendValue) TableName() string {
	return "extend_value"
}
