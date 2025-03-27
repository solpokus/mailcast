package services

import (
	"fmt"
	"mailcast/repository"
	"mailcast/util"
	"time"
)

func ProcessMsgs(body string) {
	phoneList := util.ExtractPhoneList(body)
	fmt.Printf("Extracted Phone List: %v\n", phoneList)
	isScheduleChanged := util.IsScheduleChanged(body)

	for _, phoneInfo := range phoneList {
		fmt.Printf("Phone: %s, Name: %s\n", phoneInfo.Phone, phoneInfo.Name)
		schedules := util.ExtractSchedule(body)
		segmentDetails := util.FormatSegments(schedules)

		var prevArrivalDateTime time.Time

		for i, schedule := range schedules {
			msg, location, arrivalLocation := PrepareMessageAndLocations(schedule, phoneInfo.Name)
			if location == nil || arrivalLocation == nil {
				continue // Skip if locations are not available
			}

			departDateTimeInLocation := schedule.DepartDateTime.In(location)
			arriveDateTimeInLocation := schedule.ArriveDateTime.In(arrivalLocation)

			if isScheduleChanged {
				scheduledAt := departDateTimeInLocation.Add(-24 * time.Hour)
				sendWaMessage(i, phoneInfo.Phone, msg, util.IMAGE_CHANGE, scheduledAt)
			} else {
				// HandleRegularSchedule(i, phoneInfo.Phone, msg, departDateTimeInLocation, prevArrivalDateTime, segmentDetails)
				HandleRegularScheduleAsynq(i, phoneInfo.Phone, msg, departDateTimeInLocation, prevArrivalDateTime, segmentDetails)
			}

			prevArrivalDateTime = arriveDateTimeInLocation
		}
	}
}

// Helper to prepare message, departure, and arrival locations
func PrepareMessageAndLocations(schedule util.FlightSchedule, name string) (string, *time.Location, *time.Location) {
	departFormatted := schedule.DepartDateTime.Format(util.DATE_LAYOUT)
	arriveFormatted := schedule.ArriveDateTime.Format(util.DATE_LAYOUT)

	airlineRepo, err := repository.GetAirlinesByCode(schedule.FlightNo[:2])
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
		airline,
		schedule.FlightNo,
		schedule.From, fmt.Sprintf("%s | %s", fromAirport.AirportName, fromAirport.CityName),
		schedule.To, fmt.Sprintf("%s | %s", toAirport.AirportName, toAirport.CityName),
		departFormatted,
		arriveFormatted,
	)

	return msg, location, arrivalLocation
}

// Helper to handle regular schedule notification logic
func HandleRegularSchedule(i int, phone string, msg string, departDateTime time.Time, prevArrivalDateTime time.Time, segmentDetails string) {
	if i == 0 {
		msgWithSegment := fmt.Sprintf(util.MSG_TEMPLATE_1ST, msg, segmentDetails)
		sendWaMessage(i, phone, msgWithSegment, util.IMAGE_WA_NOTIF, time.Now())

		// Send additional main ad message 24 hours before departure if needed
		scheduledAt := departDateTime.Add(-24 * time.Hour)
		sendWaMessage(i, phone, msg, util.IMAGES_ADS_MAIN, scheduledAt)
	} else {
		image := util.IMAGES_ADS_MAIN
		if departDateTime.Before(prevArrivalDateTime.Add(12 * time.Hour)) {
			image = util.IMAGE_TRANSFER
		}
		scheduledAt := departDateTime.Add(-24 * time.Hour)
		sendWaMessage(i, phone, msg, image, scheduledAt)
	}
}

// Helper to handle regular schedule notification logic to asynqmon
func HandleRegularScheduleAsynq(i int, phone string, msg string, departDateTime time.Time, prevArrivalDateTime time.Time, segmentDetails string) {
	// Date now to validate past flight
	dateTimeNow := time.Now()
	fmt.Println("dateTimeNow ", dateTimeNow)

	// Check validate past date
	pastDate := false

	fmt.Println("departDateTime , prevArrivalDateTime ", departDateTime, prevArrivalDateTime)
	if departDateTime.Before(dateTimeNow) {
		fmt.Println("departDateTime and arriveDateTimeInLocation before now")
		pastDate = true
	}

	if i == 0 {
		msgWithSegment := fmt.Sprintf(util.MSG_TEMPLATE_1ST, msg, segmentDetails)
		sendMessageToAsynq(i, phone, msgWithSegment, util.IMAGE_WA_NOTIF, time.Now())

		if !pastDate {
			// Send additional main ad message 24 hours before departure if needed
			scheduledAt := departDateTime.Add(-24 * time.Hour)
			sendMessageToAsynq(i, phone, msg, util.IMAGES_ADS_MAIN, scheduledAt)
		}
	} else {
		if !pastDate {
			image := util.IMAGES_ADS_MAIN
			// if departDateTime.Before(prevArrivalDateTime.Add(12 * time.Hour)) {
			// 	image = util.IMAGE_TRANSFER
			// }
			scheduledAt := departDateTime.Add(-24 * time.Hour)
			sendMessageToAsynq(i, phone, msg, image, scheduledAt)
		}
	}
}
