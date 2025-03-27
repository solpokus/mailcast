package repository

import (
	"log"
	"mailcast/database"
	"mailcast/models"
)

// Inserts a new log mail into the database
func InsertNewLogMail(strSubject string, strBody string, strFrom string) {

	// Insert a new record
	newMessage := models.LogMail{
		MsgSubject: strSubject,
		MsgBody:    strBody,
		MsgFrom:    strFrom,
		// ScheduleAt: scheduleAt,
		// ScheduleAt: time.Now().Add(24 * time.Hour), // Schedule for tomorrow
	}

	// Save to database
	// result := database.DB.Debug().Create(&newMessage)
	result := database.DB.Create(&newMessage)
	if result.Error != nil {
		log.Fatalf("Error inserting log_mail : %v", result.Error)
	}

	// Print inserted record ID
	log.Println("✅ Message inserted with ID:", newMessage.ID)
}

// func InsertNewLogMail(logMail *models.LogMail) error {
// 	return database.DB.Create(logMail).Error
// }
