package models

import (
	"time"
)

// LogMail represents the table structure in GORM
type TaskSchedule struct {
	ID          uint   `gorm:"primaryKey"` // serial4 maps to uint (GORM handles auto-increment)
	IdAsynq     string `gorm:"column:id_asynq"`
	Type        string `gorm:"column:type"`
	State       string `gorm:"column:state"`
	QueueName   string `gorm:"column:queue"`
	Retry       string `gorm:"column:retry"`
	Payload     string `gorm:"column:payload"`
	Completed   string `gorm:"column:completed"`
	Ttl         string `gorm:"column:ttl"`
	Status      string `gorm:"column:status"`
	ProviderPnr string `gorm:"column:provider_pnr"`
	PsgName     string `gorm:"column:psg_name"`
	PsgPhone    string `gorm:"column:psg_phone"`
	SegNo       int    `gorm:"column:segment_no"`

	CreatedDate time.Time `gorm:"column:created_date;autoCreateTime"`
	// UpdatedDate time.Time `gorm:"column:updated_date;autoUpdateTime"`
}

// TableName overrides the default table name
func (TaskSchedule) TableName() string {
	return "new.t_task_schedule"
}
