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
				"caption": msg,
			},
		},
		// "scheduledAt": scheduledAt.Format(time.RFC3339),
		"options": map[string]interface{}{
			"isSupabase": true,
		},
	}

	// Send the POST request with headers
	resp, err := clientResty.R().
		SetHeaders(map[string]string{
			// "Authorization": "defaultDS-49434e96f251d2ff",
			// "x-api-key":     "23b964f4c543ccdc",
			// "jwt":           util.JWT,
			"Accept":       "application/json",
			"Content-Type": "application/json",
			"token":        configuration.CONFIG.DaisiApiToken,
		}).
		SetBody(payload).
		Post(configuration.CONFIG.DaisiApiUrl)

	if err != nil {
		log.Fatalf("Error occurred while sending message: %v", err)
	}

	log.Println("----- Start scheduledAt -----")
	log.Println("scheduledAt :", scheduledAt)
	log.Println("----- End scheduledAt -----")

	log.Println("--------- Start Message ---------")
	// log.Println("Message :", msg)
	log.Println("Request payload: ", payload)
	log.Println("Response Status:", resp.Status())
	log.Println("--------- End Message ---------")

	log.Println("Response Status:", resp.Status())
	log.Println("Response Body:", resp.String())
}
