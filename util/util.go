package util

import (
	"fmt"
	"mailcast/repository"
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
		fmt.Println("Processing line:", line)
		if strings.Contains(line, "SegNo FlightNo Class From") {
			scheduleStart = true
			fmt.Println("scheduleStart :", scheduleStart)
			continue
		}

		if scheduleStart {
			match := regex.FindStringSubmatch(line)
			fmt.Println("len(match) :", len(match))
			if len(match) == 9 {
				segNo := match[1]
				flightNo := match[2]
				class := match[3]
				from := match[4]
				to := match[5]
				departDateTimeStr := match[6]
				arriveDateTimeStr := match[7]
				status := match[8]

				fromAirport, err := repository.GetTimezoneByCode(from)
				toAirport, err := repository.GetTimezoneByCode(to)

				location, err := time.LoadLocation(fromAirport.TzName)
				arrivalLocation, err := time.LoadLocation(toAirport.TzName)

				// Parse the date strings into time.Time
				// departDateTime, err := time.Parse(DATE_LAYOUT, departDateTimeStr)
				departDateTime, err := time.ParseInLocation(DATE_LAYOUT, departDateTimeStr, location)
				if err != nil {
					fmt.Println("Error parsing depart datetime:", err)
					continue
				}
				// arriveDateTime, err := time.Parse(DATE_LAYOUT, arriveDateTimeStr)
				arriveDateTime, err := time.ParseInLocation(DATE_LAYOUT, arriveDateTimeStr, arrivalLocation)
				if err != nil {
					fmt.Println("Error parsing arrive datetime:", err)
					continue
				}

				// Populate the schedule struct
				schedule := FlightSchedule{
					SegNo:    ParseInt(segNo),
					FlightNo: flightNo,
					Class:    class,
					From:     from,
					To:       to,
					// DepartDateTime: departDateTime.In(location),
					// ArriveDateTime: arriveDateTime.In(arrivalLocation),
					DepartDateTime: departDateTime,
					ArriveDateTime: arriveDateTime,
					Status:         status,
				}
				schedules = append(schedules, schedule)
				fmt.Printf("schedules: %s\n", schedules)
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

func ExtractPccId(body string) Pcc {
	re := regexp.MustCompile(`(?m)^PCC\s*:\s*(\S+)`)
	match := re.FindStringSubmatch(body)

	var pccList Pcc
	if len(match) > 1 {
		pccList := Pcc{
			PccId: match[1],
		}
		fmt.Println("PCC Code:", match[1])
		return pccList // found
	}

	fmt.Println("PCC not found")
	return pccList // not found
}

func ExtractProviderPnr(body string) Pnr {
	re := regexp.MustCompile(`(?m)^Provider PNR\s*:\s*(\S+)`)
	match := re.FindStringSubmatch(body)

	var pnrList Pnr
	if len(match) > 1 {
		pnrList := Pnr{
			ProviderPnr: match[1],
		}
		fmt.Println("Provider Pnr:", match[1])
		return pnrList // found
	}

	fmt.Println("Provider Pnr not found")
	return pnrList // not found
}

func ExtractPassengerList(body string) string {
	fmt.Println("start ExtractPassengerList")
	// Regex to capture everything under "Passenger List :" until the next empty line
	reBlock := regexp.MustCompile(`(?s)Passenger List\s*:\s*(.*?)\n\s*\n`)
	match := reBlock.FindStringSubmatch(body)
	// if len(match) < 2 {
	// 	return "Passenger list not found"
	// }

	var sb strings.Builder
	// sb.WriteString("Passenger Names:\n")
	lines := strings.Split(match[1], "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line != "" {
			sb.WriteString(line + "\n")
		}
	}

	fmt.Println("finish ExtractPassengerList")
	return sb.String()
}

func ExtractPassengerLastnameList(body string) string {
	fmt.Println("start ExtractPassengerLastnameList")
	// Regex to match lines under "Passenger List :"
	re := regexp.MustCompile(`(?m)^([A-Z]+),`)

	matches := re.FindAllStringSubmatch(body, -1)

	var results []string
	for _, match := range matches {
		results = append(results, match[1])
	}

	finalOutput := strings.Join(results, ", ")
	// fmt.Println(finalOutput)
	fmt.Println("finish ExtractPassengerLastnameList")
	return finalOutput
}

func HasMatchingPassenger(listDb, listParam string) bool {
	passengersA := strings.Split(listDb, ",")
	passengersB := strings.Split(listParam, ",")

	// Normalize and trim each entry
	normalize := func(list []string) []string {
		var result []string
		for _, p := range list {
			result = append(result, strings.TrimSpace(p))
		}
		return result
	}

	passengersA = normalize(passengersA)
	passengersB = normalize(passengersB)

	// Check if any passenger in A exists in B
	for _, a := range passengersA {
		for _, b := range passengersB {
			if a == b {
				return true
			}
		}
	}
	return false
}

func FormatPhoneList(phoneList []PhoneInfo) string {
	var builder strings.Builder

	for _, p := range phoneList {
		line := fmt.Sprintf("%s - %s\n", p.Phone, p.Name)
		builder.WriteString(line)
	}

	return builder.String()
}
