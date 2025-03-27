package models

import (
	"time"
)

// LogMail represents the table structure in GORM
type LogMail struct {
	ID         uint   `gorm:"primaryKey"` // serial4 maps to uint (GORM handles auto-increment)
	MsgSubject string `gorm:"column:msg_subject"`
	MsgBody    string `gorm:"column:msg_body"`
	MsgFrom    string `gorm:"column:msg_from"`
	// ScheduleAt  time.Time `gorm:"column:schedule_at"`
	CreatedDate time.Time `gorm:"column:created_date;autoCreateTime"`
	// UpdatedDate time.Time `gorm:"column:updated_date;autoUpdateTime"`
}

// TableName overrides the default table name
func (LogMail) TableName() string {
	return "new.t_log_mail"
}
