package models

// Represents the table structure in GORM
type RedisTaskDetail struct {
	ID          uint   `gorm:"primaryKey"` // serial4 maps to uint (GORM handles auto-increment)
	IdRedisTask uint   `gorm:"column:task_redis_id"`
	Key         string `gorm:"column:key"`
	Value       string `gorm:"column:value"`
}

// TableName overrides the default table name
func (RedisTaskDetail) TableName() string {
	return "new.t_redis_task_detail"
}
