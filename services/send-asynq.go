package services

import (
	"log"
	"mailcast/configuration"
	"mailcast/tasks"
	"time"

	"github.com/hibiken/asynq"
)

func sendMessageToAsynq(id int, to, msg, imgUrl string, scheduledAt time.Time, providerPnr string, psgName string, segNo int, status string, agentCredsId uint) {
	log.Printf("Sending message to asynq... ID: %d\n", id)
	log.Printf("To... : %d\n", to)

	// Define the payload
	payload := map[string]interface{}{
		"sender": configuration.CONFIG.DaisiApiSenderName,
		"phones": to,
		"messages": []map[string]interface{}{
			{
				"image": map[string]interface{}{
					"url": imgUrl,
				},
				"caption": msg,
			},
		},
		// "scheduledAt": scheduledAt.Format(time.RFC3339),
		"options": map[string]interface{}{
			"isSupabase": true,
		},
	}

	log.Println("----- Start scheduledAt -----")
	log.Println("scheduledAt :", scheduledAt)
	log.Println("scheduledAt rfc :", scheduledAt.Format(time.RFC3339))
	log.Println("----- End scheduledAt -----")

	log.Println("--------- Start Message ---------")
	// log.Println("Request payload: ", payload)
	log.Println("--------- End Message ---------")

	stcTask, err := sendToAsync(payload, to, scheduledAt, agentCredsId)
	if err != nil {
		log.Fatalf("could not send to asynq: %v", err)
	}

	InsertNewTableTaskSchedule(stcTask.ID, stcTask.Type, stcTask.Queue, string(stcTask.Payload), providerPnr, psgName, to, segNo, status)

}

func sendToAsync(payloadRes map[string]interface{}, phoneNumber string, scheduleAt time.Time, agentCredsId uint) (*asynq.TaskInfo, error) {

	log.Println("################ Start Send To Asynq ################")
	log.Println("### phoneNumber - scheduleAt", phoneNumber, " - ", scheduleAt, " ####")

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: configuration.CONFIG.RedisAddr})
	defer client.Close()

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------

	task, err := tasks.NewSchedulerTask(payloadRes, phoneNumber, scheduleAt, agentCredsId)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	info, err := client.Enqueue(
		task,
		asynq.ProcessAt(scheduleAt),
		asynq.Retention(5040*time.Hour)) // 1 Month
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	log.Println("################ Finish Send To Asynq ################")

	return info, err
}
