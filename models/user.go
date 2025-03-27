package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	FirstName string `gorm:"column:first_name"`
	Email     string `gorm:"column:email;unique"`
	CreatedAt string `gorm:"column:created_date"`
	UpdatedAt string `gorm:"column:updated_date"`
}

// TableName overrides the default table name
func (User) TableName() string {
	return "new.m_user"
}
