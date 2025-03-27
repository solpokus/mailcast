package repository

import (
	"log"
	"mailcast/database"
	"mailcast/models"
)

// GetTimezoneByCode finds a timezone by code
func GetTimezoneByCode(code string) (models.Timezonelist, error) {
	var timezones models.Timezonelist
	// result := database.DB.Debug().Unscoped().Where("airportcode = ?", code).First(&timezones)
	result := database.DB.Unscoped().Where("airportcode = ?", code).First(&timezones)

	if result.Error != nil {
		log.Println("❌ Error fetching timezones :", result.Error)
		return timezones, result.Error
	} else {
		// log.Println("✅ Timezones found using log : ", timezones.AirportName)
		return timezones, result.Error
	}
}
