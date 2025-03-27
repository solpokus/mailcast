package repository

import (
	"log"
	"mailcast/database"
	"mailcast/models"
)

// GetAirlinesByCode finds a airlines by code
func GetAirlinesByCode(code string) (models.Airlines, error) {
	var airlines models.Airlines
	// result := database.DB.Debug().Unscoped().Where("code = ?", code).First(&airlines)
	result := database.DB.Unscoped().Where("code = ?", code).First(&airlines)

	if result.Error != nil {
		log.Println("❌ Error fetching airlines :", result.Error)
		return airlines, result.Error
	} else {
		// log.Println("✅ Airlines found using log : ", airlines.Name)
		return airlines, result.Error
	}
}
