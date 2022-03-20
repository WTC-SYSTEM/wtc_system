package user

import (
	"context"
	"github.com/hawkkiller/wtc_system/user_service/pkg/logging"
)

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(userStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: userStorage,
		logger:  logger,
	}, nil
}

type Service interface {
	Create(ctx context.Context, dto CreateUserDTO) (string, error)
	GetOne(ctx context.Context, uuid string) (User, error)
	Delete(ctx context.Context, uuid string) error
	// GetByEmailAndPassword(ctx context.Context, email, password string) (User, error)
	// Update(ctx context.Context, dto UpdateUserDTO) error
}

func (s service) Create(ctx context.Context, dto CreateUserDTO) (string, error) {
	return "", nil
}

func (s service) GetOne(ctx context.Context, uuid string) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) Delete(ctx context.Context, uuid string) error {
	//TODO implement me
	panic("implement me")
}
