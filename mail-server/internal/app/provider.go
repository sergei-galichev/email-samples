package app

import (
	"mail-server/internal/delivery/email_v1"
	"mail-server/internal/delivery/http_v1"
	"mail-server/internal/model"
	repositories "mail-server/internal/repository"
	"mail-server/internal/repository/mongo/mail"
	"mail-server/internal/repository/mongo/subscriber"
	services "mail-server/internal/service"
	mailService "mail-server/internal/service/mail"
	"mail-server/pkg/storage/mongodb"
)

type serviceProvider struct {
	mailService services.MailService

	subscriberRepo repositories.SubscriberRepository
	mailRepo       repositories.MailRepository

	httpServer  *http_v1.HttpServer
	emailServer *email_v1.EmailServer

	storage *mongodb.StorageMongoDB

	mailChan chan model.Mail
}

func NewServiceProvider() *serviceProvider {
	storage := mongodb.NewStorage()

	return &serviceProvider{

		storage: storage,
	}
}

func (s *serviceProvider) SubscriberRepository() repositories.SubscriberRepository {
	if s.subscriberRepo == nil {
		s.subscriberRepo = subscriber.NewRepository(s.storage)
	}

	return s.subscriberRepo
}

func (s *serviceProvider) MailRepository() repositories.MailRepository {
	if s.mailRepo == nil {
		s.mailRepo = mail.NewRepository(s.storage)
	}
	return s.mailRepo
}

func (s *serviceProvider) MailService() services.MailService {
	if s.mailService == nil {
		s.mailService = mailService.NewService(s.SubscriberRepository(), s.MailRepository())
	}
	return s.mailService
}

func (s *serviceProvider) Server() *http_v1.HttpServer {
	if s.httpServer == nil {
		s.httpServer = http_v1.NewServer(s.MailService(), s.mailChan)
	}

	return s.httpServer
}
