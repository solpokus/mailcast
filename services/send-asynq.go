package services

import (
	"log"
	"mailcast/configuration"
	"mailcast/tasks"
	"time"

	"github.com/hibiken/asynq"
)

func sendMessageToAsynq(id int, to, msg, imgUrl string, scheduledAt time.Time) {
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
	log.Println("Request payload: ", payload)
	log.Println("--------- End Message ---------")

	// sendToAsync(payload, to, scheduledAt)

}

func sendToAsync(payloadRes map[string]interface{}, phoneNumber string, scheduleAt time.Time) {

	log.Println("################ Start Send To Asynq ################")
	log.Println("### phoneNumber - scheduleAt", phoneNumber, " - ", scheduleAt, " ####")

	client := asynq.NewClient(asynq.RedisClientOpt{Addr: configuration.CONFIG.RedisAddr})
	defer client.Close()

	// ------------------------------------------------------------
	// Example 2: Schedule task to be processed in the future.
	//            Use ProcessIn or ProcessAt option.
	// ------------------------------------------------------------

	task, err := tasks.NewSchedulerTask(payloadRes, phoneNumber, scheduleAt)
	if err != nil {
		log.Fatalf("could not create task: %v", err)
	}

	info, err := client.Enqueue(task, asynq.ProcessAt(scheduleAt))
	if err != nil {
		log.Fatalf("could not enqueue task: %v", err)
	}
	log.Printf("enqueued task: id=%s queue=%s", info.ID, info.Queue)

	log.Println("################ Finish Send To Asynq ################")
}
