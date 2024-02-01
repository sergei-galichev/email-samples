package subscriber

import (
	repositories "mail-server/internal/repository"
	"mail-server/pkg/storage/mongodb"
)

// make sure that repository implements SubscriberRepository
var _ repositories.SubscriberRepository = (*repository)(nil)

type repository struct {
	storage *mongodb.StorageMongoDB
}

func NewRepository(storage *mongodb.StorageMongoDB) *repository {
	return &repository{
		storage: storage,
	}
}
