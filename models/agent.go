package models

// Agent represents the table structure in GORM
type Agent struct {
	ID      uint   `gorm:"primaryKey"` // serial4 maps to uint (GORM handles auto-increment)
	Name    string `gorm:"column:name"`
	PicName string `gorm:"column:pic_name"`
	Email   string `gorm:"column:email"`
	Phone   string `gorm:"column:phone"`
	Country string `gorm:"column:country"`

	// CreatedDate time.Time `gorm:"column:created_date;autoCreateTime"`
	// UpdatedDate time.Time `gorm:"column:updated_date;autoUpdateTime"`
}

// TableName overrides the default table name
func (Agent) TableName() string {
	return "new.m_agent"
}
