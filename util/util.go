package util

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Utility to convert string to int
func ParseInt(value string) int {
	parsedValue, _ := strconv.Atoi(value)
	return parsedValue
}

func ExtractPhoneList(body string) []PhoneInfo {
	// Updated regex pattern to match both formats
	re := regexp.MustCompile(`(?m)^(\d{10,14})(?:/EN-|-)\d+([^\n]+)$`)
	matches := re.FindAllStringSubmatch(body, -1)

	var phoneList []PhoneInfo
	for _, match := range matches {
		if len(match) == 3 {
			phoneInfo := PhoneInfo{
				Phone: match[1],                    // The phone number part
				Name:  strings.TrimSpace(match[2]), // The name part
			}
			phoneList = append(phoneList, phoneInfo)
		}
	}

	return phoneList
}

func ExtractSchedule(body string) []FlightSchedule {
	lines := strings.Split(body, "\n")
	var schedules []FlightSchedule

	// Adjusted regular expression to allow for more flexible spacing between columns
	regex := regexp.MustCompile(`^\s*(\d+)\s+(\S+)\s+([A-Z])\s+([A-Z]{3,4})\s+([A-Z]{3,4})\s+(\d{2}\s+\w+\s+\d{4}\s+\d{2}:\d{2})\s+(\d{2}\s+\w+\s+\d{4}\s+\d{2}:\d{2})\s+(\S+)$`)

	scheduleStart := false

	for _, line := range lines {
		line = strings.TrimSpace(line)
		// fmt.Println("Processing line:", line)
		if strings.Contains(line, "SegNo FlightNo Class From  To") {
			scheduleStart = true
			continue
		}

		if scheduleStart {
			match := regex.FindStringSubmatch(line)
			if len(match) == 9 {
				segNo := match[1]
				flightNo := match[2]
				class := match[3]
				from := match[4]
				to := match[5]
				departDateTimeStr := match[6]
				arriveDateTimeStr := match[7]
				status := match[8]

				// Parse the date strings into time.Time
				departDateTime, err := time.Parse(DATE_LAYOUT, departDateTimeStr)
				if err != nil {
					fmt.Println("Error parsing depart datetime:", err)
					continue
				}
				arriveDateTime, err := time.Parse(DATE_LAYOUT, arriveDateTimeStr)
				if err != nil {
					fmt.Println("Error parsing arrive datetime:", err)
					continue
				}

				// Populate the schedule struct
				schedule := FlightSchedule{
					SegNo:          ParseInt(segNo),
					FlightNo:       flightNo,
					Class:          class,
					From:           from,
					To:             to,
					DepartDateTime: departDateTime,
					ArriveDateTime: arriveDateTime,
					Status:         status,
				}
				schedules = append(schedules, schedule)
			}
		}
	}
	return schedules
}

func FormatSegments(segments []FlightSchedule) string {
	var result string
	for _, segment := range segments {
		result += fmt.Sprintf("%d     %s    %s     %s   %s   %s %s %s\n",
			segment.SegNo, segment.FlightNo, segment.Class, segment.From, segment.To,
			segment.DepartDateTime.Format(DATE_LAYOUT), segment.ArriveDateTime.Format(DATE_LAYOUT), segment.Status)
	}
	return result
}

func IsScheduleChanged(body string) bool {
	return strings.Contains(strings.ToLower(body), "schedule change")
}
