package services

import (
	"encoding/base64"
	"fmt"
	"log"

	"google.golang.org/api/gmail/v1"
)

// Fetch and parse unread emails
func fetchEmails(srv *gmail.Service) {
	user := "me"
	query := "label:INBOX" // Filter by inbox label
	req := srv.Users.Messages.List(user)

	log.Println("Logged in")

	msgs, err := req.Q(query).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve messages: %v", err)
	}

	if len(msgs.Messages) == 0 {
		fmt.Println("No new emails found.")
		return
	}

	for _, m := range msgs.Messages {
		msg, err := srv.Users.Messages.Get(user, m.Id).Do()
		if err != nil {
			log.Fatalf("Unable to get message: %v", err)
		}

		fmt.Println("----------------------------------------------------")
		fmt.Println("========= START =========")
		fmt.Printf("Email ID: %s\n", msg.Id)
		for _, header := range msg.Payload.Headers {
			if header.Name == "From" {
				fmt.Printf("From: %s\n", header.Value)
			}
			if header.Name == "Subject" {
				fmt.Printf("Subject: %s\n", header.Value)
			}
		}

		// Decode the email body
		body := parseEmailBody(msg)
		fmt.Printf("Body: %s\n", body)
		fmt.Println("========= END =========")
	}
}

// Parse email body content
func parseEmailBody(message *gmail.Message) string {
	if message.Payload.Body.Data != "" {
		data, _ := base64.URLEncoding.DecodeString(message.Payload.Body.Data)
		return string(data)
	}

	for _, part := range message.Payload.Parts {
		if part.MimeType == "text/plain" {
			data, _ := base64.URLEncoding.DecodeString(part.Body.Data)
			return string(data)
		}
	}
	return "No text content found"
}

func parser() {
	fmt.Println("Starting Gmail Parser...")

	// Get Gmail service
	service, err := getGmailService()
	if err != nil {
		log.Fatalf("Failed to initialize Gmail API: %v", err)
	}

	// Fetch and parse emails
	fetchEmails(service)
}
