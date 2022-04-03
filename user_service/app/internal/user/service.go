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
	Create(ctx context.Context, dto CreateUserDTO) error
	GetOne(ctx context.Context, uuid string) (User, error)
	GetByEmailAndPassword(ctx context.Context, dto GetUserByEmailAndPasswordDTO) (User, error)
	Update(ctx context.Context, dto UpdateUserDTO) error
	//Delete(ctx context.Context, uuid string) error
}

func (s service) Create(ctx context.Context, dto CreateUserDTO) error {
	var user *User
	user = dto.NewUser()
	err := s.storage.Create(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (s service) GetOne(ctx context.Context, uuid string) (User, error) {
	user, err := s.storage.FindOne(ctx, uuid)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (s service) GetByEmailAndPassword(ctx context.Context, dto GetUserByEmailAndPasswordDTO) (User, error) {
	//TODO implement me
	panic("implement me")
}

func (s service) Update(ctx context.Context, dto UpdateUserDTO) error {
	//TODO implement me
	panic("implement me")
}
