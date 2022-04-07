package user

import (
	"context"
	"github.com/hawkkiller/wtc_system/user_service/internal/apperror"
	"github.com/hawkkiller/wtc_system/user_service/pkg/logging"
	"golang.org/x/crypto/bcrypt"
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
	Create(ctx context.Context, dto CreateUserHashedDTO) error
	GetOne(ctx context.Context, uuid string) (User, error)
	GetByEmailAndPassword(ctx context.Context, dto GetUserByEmailAndPasswordDTO) (User, error)
	Update(ctx context.Context, dto UpdateUserDTO) error
	//Delete(ctx context.Context, uuid string) error
}

func (s service) Create(ctx context.Context, dto CreateUserHashedDTO) error {
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
	user, err := s.storage.FindByEmail(ctx, dto.Email)
	if err != nil {
		return User{}, apperror.NewAppError("User or password is incorrect", "WTC-000005", "pass right email")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return User{}, apperror.NewAppError("User or password is incorrect", "WTC-000005", "pass right password")
	}

	return user, nil
}

func (s service) Update(ctx context.Context, dto UpdateUserDTO) error {
	//TODO implement me
	panic("implement me")
}
