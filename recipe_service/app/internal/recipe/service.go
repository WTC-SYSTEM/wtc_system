package recipe

import (
	"context"
	"github.com/hawkkiller/wtc_system/recipe_service/pkg/logging"
)

type service struct {
	storage Storage
	logger  logging.Logger
}

type Service interface {
	Create(ctx context.Context, recipe Recipe) error
}

func (s service) Create(ctx context.Context, recipe Recipe) error {
	_, err := s.storage.Create(ctx, recipe)
	if err != nil {
		return err
	}
	return nil
}

func NewService(recipeStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: recipeStorage,
		logger:  logger,
	}, nil
}
