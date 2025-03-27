package services

import (
	"fmt"
	"log"
	"mailcast/configuration"
	"time"

	"github.com/go-resty/resty/v2"
)

var clientResty = resty.New()

func sendWaMessage(id int, to, msg, imgUrl string, scheduledAt time.Time) {
	fmt.Printf("Sending WhatsApp message... ID: %d\n", id)
	log.Println("scheduledAt :", scheduledAt)
	log.Println("scheduledAt rfc :", scheduledAt.Format(time.RFC3339))

	// stringss := " /n scheduled at : "
	// addstring := stringss + scheduledAt.String()

	// Format the target number
	// toTarget := fmt.Sprintf("%s@s.whatsapp.net", to)

	// Define the payload
	payload := map[string]interface{}{
		"sender": configuration.CONFIG.DaisiApiSenderName,
		"phones": to,
		"messages": []map[string]interface{}{
			{
				"image": map[string]interface{}{
					"url": imgUrl,
				},
				// "caption": msg + addstring,
				"caption": msg,
			},
		},
		// "scheduledAt": scheduledAt.Format(time.RFC3339),
		"options": map[string]interface{}{
			"isSupabase": true,
		},
	}

	// if err != nil {
	// 	log.Fatalf("Error occurred while sending message: %v", err)
	// }

	log.Println("----- Start scheduledAt -----")
	log.Println("scheduledAt :", scheduledAt)
	log.Println("----- End scheduledAt -----")

	log.Println("--------- Start Message ---------")
	// log.Println("Message :", msg)
	log.Println("Request payload: ", payload)
	log.Println("--------- End Message ---------")

	// log.Println("Response Status:", resp.Status())
	// log.Println("Response Body:", resp.String())
}
