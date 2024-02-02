package mail

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mail-server/internal/model"
)

func (s *service) AddSubscriber(subscriber model.Subscriber) (bool, string, error) {
	return s.subscriberRepo.AddSubscriber(subscriber)
}

func (s *service) AddMail(upload model.MailUpload) (string, error) {

	return s.mailRepo.AddMail(upload)
}

func (s *service) FindSubscribers() ([]primitive.M, error) {
	return s.subscriberRepo.FindSubscribers()
}
