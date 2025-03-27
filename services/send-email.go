package services

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"mailcast/util"
	"strings"
	"time"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/emersion/go-message/mail"
	"github.com/mailgun/mailgun-go/v4"
)

const (
	EMAIL    = "ryan@timkado.id"
	PASSWORD = "danacwzdunthziui"

	apiKey = "41d360d94234707762bcca05fe721add-77316142-a21ec2ad"
	domain = "mg.daisi.app"

	SUBJECT = "26 Dec 2024 - PREFLIGHT INFO GALILEO - PR - PR2000"
	FROM    = "alert@daisi.app"
)

var TO = []string{"ryan@timkado.id", "hbinduni@yahoo.com"}

func SendMimeMessage(email util.Email) (string, error) {
	mg := mailgun.NewMailgun(domain, apiKey)

	// Use a proper MIME format, with additional headers like MIME-Version
	strMime := `MIME-Version: 1.0
Content-Type: text/plain; charset=UTF-8
Content-Transfer-Encoding: quoted-printable
Subject: %s
From: %s
To: %s
Date: %s

%s`

	// Ensure that your MIME string is properly formatted
	contentMime := fmt.Sprintf(strMime, email.Subject, email.From, strings.Join(email.To, ","), time.Now().Format(time.RFC1123Z), email.Text)

	// Create the MIME message using Mailgun
	m := mailgun.NewMIMEMessage(io.NopCloser(strings.NewReader(contentMime)), strings.Join(email.To, ","))

	// Context with timeout to avoid long waits
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	// Send the email via Mailgun
	_, id, err := mg.Send(ctx, m)
	return id, err
}

func sendEmail(content string) {

	email := util.Email{
		From:    FROM,
		To:      TO,
		Subject: SUBJECT,
		Text:    content,
	}

	id, err := SendMimeMessage(email)
	if err != nil {
		fmt.Println("Failed to send email:", err)
	} else {
		fmt.Printf("email sent successfully! %s\n", id)
	}
}

func checkEmailAndStart() {
	email := EMAIL
	password := PASSWORD

	c, err := client.DialTLS("imap.gmail.com:993", &tls.Config{})
	if err != nil {
		log.Fatal(err)
	}
	log.Println("Connected to the server")

	if err := c.Login(email, password); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged in")

	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Selected mailbox: %s with %d messages\n", mbox.Name, mbox.Messages)

	from := uint32(1)
	to := mbox.Messages
	if mbox.Messages > 5 {
		from = mbox.Messages - 5 + 1
	}

	seqSet := new(imap.SeqSet)
	seqSet.AddRange(from, to)

	messages := make(chan *imap.Message, 10)
	section := &imap.BodySectionName{}
	items := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}

	go func() {
		if err := c.Fetch(seqSet, items, messages); err != nil {
			log.Fatal(err)
		}
	}()

	processEmails(messages, section, c)

	if err := c.Logout(); err != nil {
		log.Fatal(err)
	}
	log.Println("Logged out")
}

func processEmails(messages chan *imap.Message, section *imap.BodySectionName, c *client.Client) {
	for msg := range messages {
		if msg.Envelope != nil {
			subject := msg.Envelope.Subject

			processSubject(subject, msg, section, c)
		}
	}
}

func processSubject(subject string, msg *imap.Message, section *imap.BodySectionName, c *client.Client) {
	if strings.Contains(strings.ToLower(subject), strings.ToLower("PREFLIGHT INFO GALILEO")) {
		fmt.Printf("Subject: %s\n", subject)
		fmt.Printf("date: %s", msg.Envelope.Date)

		for _, recipient := range msg.Envelope.To {
			fmt.Printf("To: %s\n", recipient.MailboxName+"@"+recipient.HostName)
		}

		if msg.Body != nil {
			shouldReturn := processBodyMsg(msg, section, c)
			if shouldReturn {
				return
			}
		}
		fmt.Println("------------- end ------------------")
	}
}

func processBodyMsg(msg *imap.Message, section *imap.BodySectionName, c *client.Client) bool {
	r := msg.GetBody(section)
	if r == nil {
		log.Println("Server didn't return message body")
		return true
	}

	// fmt.Println("body:", r)

	mr, err := mail.CreateReader(r)
	if err != nil {
		fmt.Println("Error creating mail reader:", err)
		log.Fatal(err)
	}

	for {
		part, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			fmt.Println("Error reading message part:", err)
			log.Fatal(err)
		}

		switch h := part.Header.(type) {
		case *mail.InlineHeader:
			processEmailText(h, part, c, msg, util.DELETE_EMAIL_AFTER_PROCESS)
		case *mail.AttachmentHeader:

			continue
		}
	}
	return false
}

func processEmailText(h *mail.InlineHeader, part *mail.Part, c *client.Client, msg *imap.Message, delete bool) {
	mediaType, _, _ := h.ContentType()
	if mediaType == "text/plain" || mediaType == "text/html" {
		b, _ := io.ReadAll(part.Body)
		body := strings.TrimSpace(string(b))
		body = strings.ReplaceAll(body, "\u00A0", " ")

		ProcessMsgs(body)

		if delete {

			fmt.Println("------------- before delete email ------------------")

			deleteEmail(c, msg.SeqNum)
			fmt.Println("------------- after delete email ------------------")
		}
	}
}

// Mark the email for deletion
func deleteEmail(c *client.Client, seqNum uint32) {
	seqSet := new(imap.SeqSet)
	seqSet.AddNum(seqNum)

	item := imap.FormatFlagsOp(imap.AddFlags, true)
	flags := []interface{}{imap.DeletedFlag}

	if err := c.Store(seqSet, item, flags, nil); err != nil {
		log.Fatal(err)
	}

	// Permanently delete the email
	if err := c.Expunge(nil); err != nil {
		log.Fatal(err)
	}
	log.Printf("Message %d marked for deletion and expunged\n", seqNum)
}
