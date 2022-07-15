package user

import (
	"context"
	"github.com/WTC-SYSTEM/apperror"
	"github.com/WTC-SYSTEM/logging"
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
	Create(ctx context.Context, dto CreateUserHashedDTO) (string, error)
	GetOne(ctx context.Context, uuid string) (User, error)
	GetByEmailAndPassword(ctx context.Context, dto GetUserByEmailAndPasswordDTO) (User, error)
	Update(ctx context.Context, dto UpdateUserDTO) error
	//Delete(ctx context.Context, uuid string) error
}

func (s service) Create(ctx context.Context, dto CreateUserHashedDTO) (string, error) {
	var user *User
	user = dto.NewUser()
	UUID, err := s.storage.Create(ctx, user)
	if err != nil {
		return "", err
	}
	return UUID, nil
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
		return User{}, apperror.ErrNotFound
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.Password))
	if err != nil {
		return User{}, apperror.NewAppError("User or password is incorrect", "WTC-000005", "pass right data")
	}

	return user, nil
}

func (s service) Update(ctx context.Context, dto UpdateUserDTO) error {
	// we need to get user
	user, err := s.storage.FindOne(ctx, dto.ID)
	if err != nil {
		return err
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(dto.OldPassword)); err != nil {
		return apperror.NewAppError("User or password is incorrect", "WTC-000005", "pass right data")
	}
	newUser, err := dto.NewUser(&user)
	if err != nil {
		return err
	}

	err = s.storage.Update(ctx, *newUser)

	if err != nil {
		return err
	}
	return nil
}
