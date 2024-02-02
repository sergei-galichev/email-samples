package http_v1

import (
	"fmt"
	"log"
	"mail-server/internal/model"
	"mail-server/pkg/tools"
	"net/http"
	"os"
	"time"
)

func (s *HttpServer) Home() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := tools.HTMLRender(w, req, nil)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (s *HttpServer) GetSubscriber() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var subs model.Subscriber
		subscriber, err := tools.ReadForm(req, subs)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		ok, msg, err := s.mailService.AddSubscriber(subscriber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		switch ok {
		case msg == "":
			err = tools.JSONWriter(w, "", http.StatusOK)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		case msg != "":
			err = tools.JSONWriter(w, msg, http.StatusOK)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}

	}
}

func (s *HttpServer) SendMail() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var mailUpload model.MailUpload
		upload, err := tools.ReadMultiForm(w, req, mailUpload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		msg, err := s.mailService.AddMail(upload)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println(msg)
		log.Println("........ Preparing to send mail to subscribers ........ ")
		time.Sleep(100 * time.Millisecond)
		log.Println("........ Accessing the subscribers Database ........... ")

		res, err := s.mailService.FindSubscribers()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		for _, item := range res {
			subEmail := item["mail"].(string)
			firstName := item["first_name"].(string)
			lastName := item["last_name"].(string)

			subName := fmt.Sprintf("%s %s", firstName, lastName)
			mail := model.Mail{
				Source:      os.Getenv("GMAIL_ACC"),
				Destination: subEmail,
				Name:        subName,
				Message:     upload.DocxContent,
				Subject:     upload.DocxName,
			}

			s.mailChan <- mail
		}

		err = tools.JSONWriter(w, fmt.Sprintf("Mail sent %v subscribers", len(res)), http.StatusOK)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	}
}
