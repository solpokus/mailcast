package repository

import (
	"log"
	"mailcast/database"
	"mailcast/models"
)

// Inserts a new task_schedule into the database
func InsertNewTaskSchedule(idAsynq string, typeStr string, queueName string, payload string, providerPnr string,
	psgName string, psgPhone string, segNo int, status string) {

	// Insert a new record
	newMessage := models.TaskSchedule{
		IdAsynq:     idAsynq,
		Type:        typeStr,
		QueueName:   queueName,
		Payload:     payload,
		ProviderPnr: providerPnr,
		PsgName:     psgName,
		PsgPhone:    psgPhone,
		SegNo:       segNo,
		Status:      status,
	}

	// Save to database
	// result := database.DB.Debug().Create(&newMessage)
	result := database.DB.Create(&newMessage)
	if result.Error != nil {
		log.Fatalf("Error inserting task_schedule : %v", result.Error)
	}

	// Print inserted record ID
	log.Println("✅ Message inserted to task_schedule with ID:", newMessage.ID)
}

// GetTaskByIdAsynq finds a TaskSchedule by id_asynq
func GetTaskByProviderPnrAndPsgName(providerPnr string, psgName string, segNo int) ([]models.TaskSchedule, error) {
	var taskSchedule []models.TaskSchedule
	result := database.DB.Unscoped().
		Where("provider_pnr = ?", providerPnr).
		Where("segment_no = ?", segNo).
		Where("status = ?", "HK").
		Where("psg_name LIKE ?", "%"+psgName+"%").Find(&taskSchedule)

	if result.Error != nil {
		log.Println("❌ Error fetching TaskSchedule :", result.Error)
		return taskSchedule, result.Error
	} else {
		return taskSchedule, result.Error
	}
}

// V1
func GetTaskByProviderPnrAndPsgNameV1(providerPnr string, psgName string) ([]models.TaskSchedule, error) {
	var taskSchedule []models.TaskSchedule
	result := database.DB.Unscoped().
		Where("provider_pnr = ?", providerPnr).
		Where("psg_name LIKE ?", "%"+psgName+"%").Find(&taskSchedule)

	if result.Error != nil {
		log.Println("❌ Error fetching GetTaskByProviderPnrAndPsgNameV1 :", result.Error)
		return taskSchedule, result.Error
	} else {
		return taskSchedule, result.Error
	}
}

// Find by providerPnr
func GetTaskByProviderPnr(providerPnr string) ([]models.TaskSchedule, error) {
	var taskSchedule []models.TaskSchedule
	result := database.DB.Unscoped().
		Where("provider_pnr = ?", providerPnr).Find(&taskSchedule)

	if result.Error != nil {
		log.Println("❌ Error fetching TaskSchedule :", result.Error)
		return taskSchedule, result.Error
	} else {
		return taskSchedule, result.Error
	}
}
