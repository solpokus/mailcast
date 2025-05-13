package repository

import (
	"log"
	"mailcast/database"
	"mailcast/models"
)

func GetAgentCredentialsByPccAndMailbox(pccCode string, mailbox string) (models.AgentCredential, error) {
	var agentCredential models.AgentCredential
	result := database.DB.Unscoped().
		Where("pcc_code = ?", pccCode).
		Where("daisi_mailbox = ?", mailbox).Find(&agentCredential)

	if result.Error != nil {
		log.Println("❌ Error fetching GetAgentCredentialsByPccAndMailbox : ", result.Error)
		return agentCredential, result.Error
	} else {
		return agentCredential, result.Error
	}
}
