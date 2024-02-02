package mail

import (
	repositories "mail-server/internal/repository"
	services "mail-server/internal/service"
)

var _ services.MailService = (*service)(nil)

type service struct {
	subscriberRepo repositories.SubscriberRepository
	mailRepo       repositories.MailRepository
}

func NewService(
	subscriberRepo repositories.SubscriberRepository,
	mailRepo repositories.MailRepository,
) *service {
	return &service{
		subscriberRepo: subscriberRepo,
		mailRepo:       mailRepo,
	}
}
