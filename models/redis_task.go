package models

// Represents the table structure in GORM
type RedisTask struct {
	ID      uint   `gorm:"primaryKey"` // serial4 maps to uint (GORM handles auto-increment)
	IdAsynq string `gorm:"column:id_asynq"`
}

// TableName overrides the default table name
func (RedisTask) TableName() string {
	return "new.t_redis_task"
}
