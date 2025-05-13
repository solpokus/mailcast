package services

import (
	"log"
	"mailcast/repository"
)

func InsertNewTableTaskSchedule(idAsynq string, typeStr string, queueName string, payload string, providerPnr string,
	psgName string, psgPhone string, segNo int, status string) {

	log.Println(">>>>> Start InsertNewTableTaskSchedule <<<<<")

	repository.InsertNewTaskSchedule(idAsynq, typeStr, queueName, payload, providerPnr, psgName, psgPhone, segNo, status)

	log.Println(">>>>> Finish InsertNewTableTaskSchedule <<<<<")
}
