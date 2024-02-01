package repository

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"mail-server/internal/model"
)

type SubscriberRepository interface {
	AddSubscriber(subscriber model.Subscriber) (bool, string, error)
	FindSubscribers() ([]primitive.M, error)
}

type MailRepository interface {
	AddMail(upload model.MailUpload) (string, error)
}
