package models

import "gorm.io/gorm"

type Airlines struct {
	gorm.Model
	Code string `gorm:"column:code;unique"`
	Name string `gorm:"column:name;unique"`
}

// TableName overrides the default table name
func (Airlines) TableName() string {
	return "new.m_airlines"
}
