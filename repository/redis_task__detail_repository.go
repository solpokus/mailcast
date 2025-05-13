package repository

import (
	"log"
	"mailcast/database"
	"mailcast/models"
)

// Inserts a new redis_task_detail into the database
func InsertNewRedisTaskDetail(taskRedisId uint, key string, value string) uint {

	// Insert a new record
	newMessage := models.RedisTaskDetail{
		IdRedisTask: taskRedisId,
		Key:         key,
		Value:       value,
	}

	// Save to database
	// result := database.DB.Debug().Create(&newMessage)
	result := database.DB.Create(&newMessage)
	if result.Error != nil {
		log.Fatalf("Error inserting redis_task_detail : %v", result.Error)
	}

	// Print inserted record ID
	log.Println("✅ Message inserted to redis_task_detail with ID:", newMessage.ID)

	return newMessage.ID
}
