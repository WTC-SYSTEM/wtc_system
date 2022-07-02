package photo

import (
	"context"
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
)

type service struct {
	storage Storage
	logger  logging.Logger
}

type PhotoService interface {
	Upload(ctx context.Context) error
}

func NewService(storage Storage, logger logging.Logger) PhotoService {
	return &service{
		storage: storage,
		logger:  logger,
	}
}

func (s service) Upload(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
