package services

import (
	"fmt"
	"log"
	"mailcast/models"
	"mailcast/repository"
	"mailcast/tasks"
	"mailcast/util"
	"strings"
	"time"
)

func ProcessMsgs(body string, mailboxEmail string) {
	phoneList := util.ExtractPhoneList(body)
	fmt.Printf("Extracted Phone List: %v\n", phoneList)
	isScheduleChanged := util.IsScheduleChanged(body)
	pnrList := util.ExtractProviderPnr(body)
	passengerList := util.ExtractPassengerLastnameList(body)
	listTaskFresh := getListTaskByPnrAndLastName(pnrList.ProviderPnr, passengerList)
	pccId := util.ExtractPccId(body).PccId

	//Message
	// msgPassengerList := util.ExtractPassengerList(body)
	msgPhoneList := util.FormatPhoneList(phoneList)

	//Get Agent
	agentCredential, _ := repository.GetAgentCredentialsByPccAndMailbox(pccId, mailboxEmail)
	agentCredentialId := agentCredential.ID

	agentMaster, _ := repository.GetAgentById(agentCredential.AgentId)

	for _, phoneInfo := range phoneList {
		fmt.Printf("Phone: %s, Name: %s\n", phoneInfo.Phone, phoneInfo.Name)
		schedules := util.ExtractSchedule(body)
		segmentDetails := util.FormatSegments(schedules)

		var prevArrivalDateTime time.Time

		for i, schedule := range schedules {
			fmt.Printf("schedule: %s\n", schedule.FlightNo)
			msg, location, arrivalLocation := PrepareMessageAndLocations(schedule, phoneInfo.Name, msgPhoneList, agentMaster.Name)
			if location == nil || arrivalLocation == nil {
				continue // Skip if locations are not available
			}

			// departDateTimeInLocation := schedule.DepartDateTime.In(location)
			// arriveDateTimeInLocation := schedule.ArriveDateTime.In(arrivalLocation)

			if isScheduleChanged {
				// If Status not approved non HK
				log.Printf("Segment number : %v", schedule.SegNo, "segment status : %v", schedule.Status)
				if listTaskFresh != nil {
					log.Printf("Start schedule change : %v ", i)
					scheduleChange(i, msg, schedule, phoneInfo, pnrList.ProviderPnr, passengerList, listTaskFresh, agentCredentialId)
					log.Printf("Finish schedule change : %v ", i)
				}
			} else {
				// HandleRegularSchedule(i, phoneInfo.Phone, msg, departDateTimeInLocation, prevArrivalDateTime, segmentDetails)
				HandleRegularScheduleAsynq(i, phoneInfo.Phone, phoneInfo.Name, msg, schedule.DepartDateTime, prevArrivalDateTime,
					segmentDetails, pnrList.ProviderPnr, passengerList, schedule.SegNo, schedule.Status, msgPhoneList, agentCredentialId, agentMaster.Name)
			}

			prevArrivalDateTime = schedule.ArriveDateTime
		}
	}
}

// Helper to prepare message, departure, and arrival locations
func PrepareMessageAndLocations(schedule util.FlightSchedule, name string, msgPhoneList string, agentName string) (string, *time.Location, *time.Location) {
	departFormatted := schedule.DepartDateTime.Format(util.DATE_LAYOUT)
	arriveFormatted := schedule.ArriveDateTime.Format(util.DATE_LAYOUT)

	prefix := strings.ReplaceAll(schedule.FlightNo, "/", "")
	if len(prefix) >= 2 {
		prefix = prefix[:2]
	}

	airlineRepo, err := repository.GetAirlinesByCode(prefix)
	airline := airlineRepo.Name
	fromAirport, err := repository.GetTimezoneByCode(schedule.From)
	toAirport, err := repository.GetTimezoneByCode(schedule.To)

	location, err := time.LoadLocation(fromAirport.TzName)
	if err != nil {
		fmt.Println("Error loading departure location:", err)
		return "", nil, nil
	}

	arrivalLocation, err := time.LoadLocation(toAirport.TzName)
	if err != nil {
		fmt.Println("Error loading arrival location:", err)
		return "", nil, nil
	}

	msg := fmt.Sprintf(util.MSG_TEMPLATE,
		name,
		msgPhoneList,
		airline,
		schedule.FlightNo,
		schedule.From, fmt.Sprintf("%s | %s", fromAirport.AirportName, fromAirport.CityName),
		schedule.To, fmt.Sprintf("%s | %s", toAirport.AirportName, toAirport.CityName),
		departFormatted,
		arriveFormatted,
		agentName,
	)

	return msg, location, arrivalLocation
}

// Helper to handle regular schedule notification logic
// func HandleRegularSchedule(i int, phone string, msg string, departDateTime time.Time, prevArrivalDateTime time.Time, segmentDetails string) {
// 	if i == 0 {
// 		msgWithSegment := fmt.Sprintf(util.MSG_TEMPLATE_1ST, msg, segmentDetails)
// 		sendWaMessage(i, phone, msgWithSegment, util.IMAGE_WA_NOTIF, time.Now())

// 		// Send additional main ad message 24 hours before departure if needed
// 		scheduledAt := departDateTime.Add(-24 * time.Hour)
// 		sendWaMessage(i, phone, msg, util.IMAGES_ADS_MAIN, scheduledAt)
// 	} else {
// 		image := util.IMAGES_ADS_MAIN
// 		if departDateTime.Before(prevArrivalDateTime.Add(12 * time.Hour)) {
// 			image = util.IMAGE_TRANSFER
// 		}
// 		scheduledAt := departDateTime.Add(-24 * time.Hour)
// 		sendWaMessage(i, phone, msg, image, scheduledAt)
// 	}
// }

// Helper to handle regular schedule notification logic to asynqmon
func HandleRegularScheduleAsynq(i int, phone string, phoneName string, msg string, departDateTime time.Time, prevArrivalDateTime time.Time, segmentDetails string,
	providerPnr string, psgName string, segNo int, status string, msgPhoneList string, agentCredsId uint, travelAgentName string) {
	// Image rotate
	imagesList := []string{
		util.IMAGES_ADS_MAIN_1,
		util.IMAGES_ADS_MAIN_2,
		util.IMAGES_ADS_MAIN_3,
		util.IMAGES_ADS_MAIN_4,
		util.IMAGES_ADS_MAIN_5,
	}

	image := imagesList[i%len(imagesList)] // cycle through images
	fmt.Printf("Data %d uses image: %s\n", i+1, image)

	// Date now to validate past flight
	dateTimeNow := time.Now()

	// Check validate past date
	pastDate := false

	fmt.Println("departDateTime , prevArrivalDateTime ", departDateTime, prevArrivalDateTime)
	if departDateTime.Before(dateTimeNow) && prevArrivalDateTime.Before(dateTimeNow) {
		fmt.Println("departDateTime and arriveDateTimeInLocation before now")
		pastDate = true
	}

	if i == 0 {
		msgWithSegment := fmt.Sprintf(util.MSG_TEMPLATE_1ST, phoneName, travelAgentName, msgPhoneList, segmentDetails)
		sendMessageToAsynq(i, phone, msgWithSegment, util.IMAGE_WA_NOTIF, time.Now().Add(2*time.Second), providerPnr, psgName, segNo, status, agentCredsId)

		if !pastDate {
			// Send additional main ad message 24 hours before departure if needed
			scheduledAt := departDateTime.Add(-24 * time.Hour)
			// sendMessageToAsynq(i, phone, msg, util.IMAGES_ADS_MAIN, scheduledAt)

			if scheduledAt.Before(time.Now()) {
				sendMessageToAsynq(i, phone, msg, image, time.Now().Add(10*time.Second), providerPnr, psgName, segNo, status, agentCredsId)
			} else {
				sendMessageToAsynq(i, phone, msg, image, scheduledAt, providerPnr, psgName, segNo, status, agentCredsId)
			}
		}
	} else {
		if !pastDate {
			// image := util.IMAGES_ADS_MAIN
			// if departDateTime.Before(prevArrivalDateTime.Add(12 * time.Hour)) {
			// 	image = util.IMAGE_TRANSFER
			// }
			scheduledAt := departDateTime.Add(-24 * time.Hour)
			// sendMessageToAsynq(i, phone, msg, image, scheduledAt)

			if scheduledAt.Before(time.Now()) {
				sendMessageToAsynq(i, phone, msg, image, time.Now().Add(10*time.Second), providerPnr, psgName, segNo, status, agentCredsId)
			} else {
				sendMessageToAsynq(i, phone, msg, image, scheduledAt, providerPnr, psgName, segNo, status, agentCredsId)
			}
		}
	}
}

func scheduleChange(i int, msg string, schedule util.FlightSchedule, phoneInfo util.PhoneInfo,
	providerPnr string, psgList string, taskFresh []models.TaskSchedule, agentCredsId uint) {

	idSet := make(map[uint]bool)
	var idList []string

	// Find Existing task on postgre and redis
	if schedule.Status == "TK" {
		for _, task := range taskFresh {
			if task.SegNo == schedule.SegNo {
				// inside your loop
				if !idSet[task.ID] {
					idSet[task.ID] = true
					idList = append(idList, task.IdAsynq)
				}
			}
		}

		log.Printf("idList : %v\n", idList)

		// Delete schedule on asynq
		if idList != nil {
			for _, idAsynq := range idList {
				isDelete, err := tasks.DeleteOneTasks(idAsynq, "default")
				if err != nil {
					fmt.Println("Error DeleteOneTasks asynq :", err)
				}

				if isDelete {
					fmt.Println("Success DeleteOneTasks asynq_id :", idAsynq)
				}
			}
		}

		scheduledAt := schedule.DepartDateTime.Add(-24 * time.Hour)
		sendMessageToAsynq(i, phoneInfo.Phone, msg, util.IMAGE_CHANGE, scheduledAt, providerPnr, psgList, schedule.SegNo, schedule.Status, agentCredsId)
	}
}

func getListTaskByPnrAndLastName(providerPnr string, psgList string) []models.TaskSchedule {
	taskScheduleByPnrList, err := repository.GetTaskByProviderPnr(providerPnr)
	if err != nil {
		fmt.Println("Error taskScheduleByPnrList on postgres :", err)
	} else {

		for _, taskScheduleByPnr := range taskScheduleByPnrList {
			if util.HasMatchingPassenger(taskScheduleByPnr.PsgName, psgList) {
				taskScheduleByPnrAndPsgNameList, err := repository.GetTaskByProviderPnrAndPsgNameV1(providerPnr, psgList)
				if err != nil {
					fmt.Println("Error get task on postgres :", err)
				} else {
					return taskScheduleByPnrAndPsgNameList
				}
			}
		}
	}
	return nil
}
