package util

import "time"

type Email struct {
	From    string
	To      []string
	Subject string
	Text    string
}

type PhoneInfo struct {
	Phone string
	Name  string
}

type FlightSchedule struct {
	SegNo          int
	FlightNo       string
	Class          string
	From           string
	To             string
	DepartDateTime time.Time
	ArriveDateTime time.Time
	Status         string
}

type Schedule struct {
	DepartDateTime time.Time // Departure date and time
	ArriveDateTime time.Time // Arrival date and time
	FlightNo       string    // Flight number
	From           string    // Departure airport code
	To             string    // Arrival airport code
}
