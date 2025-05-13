package tasks

import (
	"encoding/json"
	"fmt"
	"log"
	"mailcast/configuration"
	"time"

	"github.com/hibiken/asynq"
)

// A list of task types.
const (
	TypeNotif = "type:notif"
)

type SchedulerPayload struct {
	Payload      map[string]interface{}
	Phone        string
	ScheduleAt   time.Time
	AgentCredsId uint
}

//----------------------------------------------
// Write a function NewXXXTask to create a task.
// A task consists of a type and a payload.
//----------------------------------------------

func NewSchedulerTask(payloadRes map[string]interface{}, phoneNumber string, scheduleAt time.Time, agentCredsId uint) (*asynq.Task, error) {
	payload, err := json.Marshal(SchedulerPayload{Payload: payloadRes, Phone: phoneNumber, ScheduleAt: scheduleAt, AgentCredsId: agentCredsId})
	if err != nil {
		return nil, err
	}
	return asynq.NewTask(TypeNotif, payload), nil
}

func ListScheduledTasks() error {
	queue := "default"

	allSchedule, nil := FetchAllScheduledTasks(queue)

	for _, t := range allSchedule {
		fmt.Printf("Task ID: %s\n", t.ID)
		fmt.Printf("Payload: %s\n", t.Payload)
		fmt.Printf("Type: %s\n", t.Type)
		fmt.Println("---------------------------")
	}

	return nil
}

func FetchOneTasks(taskId string, queueName string) (*asynq.TaskInfo, error) {
	// Initialize the Inspector with Redis connection options
	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: configuration.CONFIG.RedisAddr})

	// Retrieve task information
	taskInfo, err := inspector.GetTaskInfo(queueName, taskId)
	if err != nil {
		log.Fatalf("Could not retrieve task info: %v", err)
	}

	// Print task details
	fmt.Printf("Task ID: %s\n", taskInfo.ID)
	fmt.Printf("Task Type: %s\n", taskInfo.Type)
	fmt.Printf("Payload: %s\n", taskInfo.Payload)
	fmt.Printf("State: %s\n", taskInfo.State)
	fmt.Printf("Next Process At: %v\n", taskInfo.NextProcessAt)
	fmt.Printf("Result: %s\n", taskInfo.Result)
	return taskInfo, nil
}

func DeleteOneTasks(taskId string, queueName string) (bool, error) {
	// Initialize the Inspector with Redis connection options
	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: configuration.CONFIG.RedisAddr})

	// Retrieve task information
	taskInfo, err := inspector.GetTaskInfo(queueName, taskId)
	if err != nil {
		log.Println("Could not retrieve task info: %v", err)
		return false, err
	}

	log.Println("xxxxxxx Start Task deleted xxxxxxx : " + taskId)

	// Print task details
	fmt.Printf("Task ID: %s\n", taskInfo.ID)
	fmt.Printf("Task Type: %s\n", taskInfo.Type)
	// fmt.Printf("Payload: %s\n", taskInfo.Payload)
	fmt.Printf("State: %s\n", taskInfo.State)
	fmt.Printf("Next Process At: %v\n", taskInfo.NextProcessAt)

	// Delete the task by ID
	if err := inspector.DeleteTask(queueName, taskId); err != nil {
		log.Println("Failed to delete task: %v", err)
		return false, err
	}

	log.Println("✅ xxxxxxx Finish Task deleted successfully xxxxxxx : " + taskId)

	return true, nil

}

func FetchAllScheduledTasks(queue string) ([]*asynq.TaskInfo, error) {
	inspector := asynq.NewInspector(asynq.RedisClientOpt{Addr: configuration.CONFIG.RedisAddr})
	var allTasks []*asynq.TaskInfo

	pageSize := 100
	page := 0

	for {
		tasks, err := inspector.ListScheduledTasks(queue, page, pageSize)
		if err != nil {
			return nil, err
		}
		allTasks = append(allTasks, tasks...)

		if len(tasks) < pageSize {
			break // last page
		}
		page++ // move to next page
	}

	log.Printf("Total tasks fetched: %d", len(allTasks))
	return allTasks, nil
}
