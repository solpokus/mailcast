package services

import (
	"fmt"
	"log"
	"mailcast/repository"
	"mailcast/util"
	"regexp"
	"strings"
	"sync"

	"google.golang.org/api/gmail/v1"
)

func SchedulerEmail() error {
	log.Println("✅ Task SchedulerEmail executed successfully")
	err := CheckEmailOauthAndStart()
	if err != nil {
		return fmt.Errorf("Failed to to run SchedulerEmail : %w", err)
	}
	return nil
}

func CheckEmailOauthAndStart() error {

	fmt.Println("Starting Gmail Parser...")

	// Get Gmail service
	srv, err := getGmailService()
	if err != nil {
		log.Println("Failed to initialize Gmail API: %v", err)
		return fmt.Errorf("Failed to initialize Gmail API: %w", err)
	}

	log.Println("Connected to the server")

	user := "me"
	query := "label:INBOX" // Filter by inbox label
	req := srv.Users.Messages.List(user)

	// Get Email
	emailUrl := getEmail()
	log.Println("Using Email User : " + emailUrl)

	log.Println("Logged in")

	msgs, err := req.Q(query).Do()
	if err != nil {
		log.Println("Unable to retrieve messages: %v", err)
		return fmt.Errorf("Unable to retrieve messages: %w", err)
	}

	log.Printf("Fetching %d latest emails asynchronously...\n", len(msgs.Messages))

	var wg sync.WaitGroup
	// emailChan := make(chan string, len(msgs.Messages))
	emailChan := make(chan *gmail.Message, len(msgs.Messages))

	for _, msg := range msgs.Messages {
		wg.Add(1)
		go func(msgID string) {
			defer wg.Done()

			message, err := srv.Users.Messages.Get(user, msgID).Do()
			if err != nil {
				log.Printf("Failed to get message %s: %v", msgID, err)
				return
			}

			// fmt.Println("----------------------------------------------------")
			// fmt.Println("========= START =========")
			// fmt.Printf("Email ID: %s\n", msgID)

			// Send email details to channel
			// emailChan <- fmt.Sprintf("From: %s\nSubject: %s\nBody: %s\n", sender, subject, body)
			emailChan <- message
		}(msg.Id)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	close(emailChan)

	// Print emails
	for email := range emailChan {
		// fmt.Println("----------------------------------------------------")
		// fmt.Println(email)
		processEmailsV2(email, srv, emailUrl)
	}

	log.Println("Logged out")

	return nil
}

func processEmailsV2(messages *gmail.Message, service *gmail.Service, mailboxEmail string) {

	// fmt.Println("----------------------------------------------------")
	// fmt.Println(messages)
	if messages != nil {

		subject := getHeaderByName("Subject", messages)

		processSubjectV2(subject, messages, service, messages.Id, mailboxEmail)
	}
}

func processSubjectV2(subject string, msg *gmail.Message, service *gmail.Service, messageID string, mailboxEmail string) {
	if strings.Contains(strings.ToLower(subject), strings.ToLower("PREFLIGHT INFO GALILEO")) {
		// if strings.Contains(strings.ToLower(subject), strings.ToLower("PREFLIGHT INFO GALILEO - SQ - B2B3B4")) {
		fmt.Printf("Subject: %s\n", subject)
		fmt.Printf("date: %s\n", getHeaderByName("Date", msg))

		recipients := getRecipientEmails(msg)
		fmt.Println("Recipients:", recipients)

		shouldReturn := processBodyMsgV2(msg, service, messageID, mailboxEmail)
		if shouldReturn {
			fmt.Println("s:", shouldReturn)
			return
		}
	}
}

func processBodyMsgV2(msg *gmail.Message, service *gmail.Service, messageID string, mailboxEmail string) bool {
	// Decode the email body
	r := parseEmailBody(msg)

	if r == "" {
		log.Println("Server didn't return message body")
		return true
	}

	// fmt.Println("start body:", r)

	processEmailTextV2(r, msg, service, messageID, util.DELETE_EMAIL_AFTER_PROCESS, mailboxEmail)

	return false
}

func processEmailTextV2(body string, msg *gmail.Message, service *gmail.Service, messageID string, delete bool, mailboxEmail string) {

	fmt.Println("------------- start processMsgs ------------------")
	b := strings.TrimSpace(string(body))
	b = strings.ReplaceAll(b, "\u00A0", " ")

	// Uncomment later
	ProcessMsgs(body, mailboxEmail)
	// fmt.Println("start processMsgs body:", b)

	// Insert mail log to db
	insertLog(msg, body)

	if delete {

		fmt.Println("--- before delete email ---")

		// deleteEmail(c, msg.SeqNum)
		// deleteEmailPermanently(service, messageID)
		trashEmail(service, messageID)
		fmt.Println("--- after delete email ---")
	}

	fmt.Println("-------------- end processMsgs -------------------")

}

// Extract recipient email addresses
func getRecipientEmails(message *gmail.Message) []string {
	var recipients []string
	emailRegex := regexp.MustCompile(`[\w\.\-]+@[\w\.\-]+\.\w+`) // Regex to extract email addresses

	for _, header := range message.Payload.Headers {
		if header.Name == "To" { // "To" contains recipient emails
			emails := emailRegex.FindAllString(header.Value, -1) // Extract emails
			recipients = append(recipients, emails...)
		}
	}

	return recipients
}

func getHeaderByName(headerName string, messages *gmail.Message) string {
	for _, header := range messages.Payload.Headers {
		if header.Name == headerName {
			return header.Value
		}
	}
	return ""
}

func deleteEmailPermanently(service *gmail.Service, messageID string) error {
	err := service.Users.Messages.Delete("me", messageID).Do()
	if err != nil {
		return fmt.Errorf("unable to delete email: %v", err)
	}
	fmt.Println("Email permanently deleted:", messageID)
	return nil
}

func trashEmail(service *gmail.Service, messageID string) error {
	_, err := service.Users.Messages.Trash("me", messageID).Do()
	fmt.Println("messageID:", messageID)
	if err != nil {
		return fmt.Errorf("unable to move email to trash: %v", err)
	}
	fmt.Println("Email moved to Trash:", messageID)
	return nil
}

func insertLog(msg *gmail.Message, body string) {
	log.Println(">>>>> Insert Log Mail <<<<<")

	subject := getHeaderByName("Subject", msg)
	from := getHeaderByName("From", msg)

	pccList := util.ExtractPccId(body)

	repository.InsertNewLogMail(subject, body, from, pccList.PccId)
}
