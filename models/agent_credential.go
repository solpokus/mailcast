package models

import (
	"time"
)

// AgentCredential represents the table structure in GORM
type AgentCredential struct {
	ID                 uint   `gorm:"primaryKey"` // serial4 maps to uint (GORM handles auto-increment)
	PccCode            string `gorm:"column:pcc_code"`
	AgentId            uint   `gorm:"column:agent_id"`
	DaisiMailbox       string `gorm:"column:daisi_mailbox"`
	DaisiApiUrl        string `gorm:"column:daisi_api_url"`
	DaisiApiSenderName string `gorm:"column:daisi_api_sender_name"`
	DaisiApiToken      string `gorm:"column:daisi_api_token"`

	CreatedDate time.Time `gorm:"column:created_date;autoCreateTime"`
	UpdatedDate time.Time `gorm:"column:updated_date;autoUpdateTime"`
}

// TableName overrides the default table name
func (AgentCredential) TableName() string {
	return "new.m_agent_credential"
}
