package photo

import (
	"context"
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
	"github.com/WTC-SYSTEM/wtc_system/photo_service/pkg/client/aws"
)

type Storage interface {
	Create(ctx context.Context) error
}

type storage struct {
	logger logging.Logger
	awsCfg aws.Aws
}

// NewStorage
func NewStorage(l logging.Logger, awsCfg aws.Aws) Storage {
	return &storage{
		logger: l,
		awsCfg: awsCfg,
	}
}

func (s storage) Create(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}
