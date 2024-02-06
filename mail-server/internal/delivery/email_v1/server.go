package email_v1

import (
	"gopkg.in/gomail.v2"
	"log"
	"mail-server/internal/model"
	"os"
)

type EmailServer struct {
	dialer *gomail.Dialer

	mailChan    chan model.Mail
	workerCount int
}

func NewEmailServer(mailChan chan model.Mail, workerCount int) *EmailServer {
	account := os.Getenv("GMAIL_ACC")
	appPass := os.Getenv("APP_PASSWORD")

	dialer := gomail.NewDialer("smtp.gmail.com", 465, account, appPass)
	s, err := dialer.Dial()
	if err != nil {
		log.Fatal("[EmailServer] NewEmailServer: error dialing")
	}

	defer func() {
		err = s.Close()
		if err != nil {
			log.Fatal("[EmailServer] NewEmailServer: error closing SendCloser")
		}
	}()
	return &EmailServer{
		dialer:      dialer,
		mailChan:    mailChan,
		workerCount: workerCount,
	}
}

func (e *EmailServer) sendEmail(mailChan model.Mail) {
	s, err := e.dialer.Dial()
	if err != nil {
		log.Println("[EmailServer] SendMail: error dialing")
	}
	defer func() {
		err = s.Close()
		if err != nil {
			log.Fatal("[EmailServer] SendMail: error closing SendCloser")
		}
	}()

	userName := os.Getenv("USER_NAME")

	msg := gomail.NewMessage()

	msg.SetAddressHeader("From", mailChan.Source, userName)
	msg.SetHeader("To", mailChan.Destination)
	msg.SetHeader("Subject", mailChan.Subject)
	msg.SetBody("text/html", mailChan.Message)
	if err = gomail.Send(s, msg); err != nil {
		log.Println("[EmailServer] SendMail: error send email")
	}

	msg.Reset()
}

func (e *EmailServer) MailDelivery() {
	completionChan := make(chan bool, e.workerCount)

	for x := 0; x < e.workerCount; x++ {
		go func() {
			defer func() {
				completionChan <- true
			}()
			for m := range e.mailChan {
				e.sendEmail(m)
			}
		}()
	}

	log.Println("Email Server is started")

	for x := 0; x < e.workerCount; x++ {
		<-completionChan
	}
}
