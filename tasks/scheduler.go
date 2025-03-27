package tasks

import (
	"encoding/json"
	"time"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeNotif = "type:notif"
)

type SchedulerPayload struct {
	Payload    map[string]interface{}
	Phone      string
	ScheduleAt time.Time
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func NewSchedulerTask(payloadRes map[string]interface{}, phoneNumber string, scheduleAt time.Time) (*asynq.Task, error) {
	payload, err := json.Marshal(SchedulerPayload{Payload: payloadRes, Phone: phoneNumber, ScheduleAt: scheduleAt})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeNotif, payload), nil
}
