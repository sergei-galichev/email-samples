package service

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mail-server/internal/model"
)

type MailService interface {
	AddSubscriber(subscriber model.Subscriber) (bool, string, error)
	AddMail(upload model.MailUpload) (string, error)
	FindSubscribers() ([]primitive.M, error)
}
