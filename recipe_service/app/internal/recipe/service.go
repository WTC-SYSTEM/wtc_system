package recipe

import "github.com/hawkkiller/wtc_system/recipe_service/pkg/logging"

type service struct {
	storage Storage
	logger  logging.Logger
}

func NewService(recipeStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: recipeStorage,
		logger:  logger,
	}, nil
}

type Service interface {
}
