package repository

import (
	"log"
	"mailcast/database"
	"mailcast/models"
)

func GetAgentById(idAgent uint) (models.Agent, error) {
	var masterAgent models.Agent
	result := database.DB.Unscoped().
		Where("id = ?", idAgent).Find(&masterAgent)

	if result.Error != nil {
		log.Println("❌ Error fetching GetAgentById : ", result.Error)
		return masterAgent, result.Error
	} else {
		return masterAgent, result.Error
	}
}
