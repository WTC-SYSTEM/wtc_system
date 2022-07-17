package recipe

import (
	"context"
	"github.com/WTC-SYSTEM/apperror"
	"github.com/WTC-SYSTEM/logging"
)

type service struct {
	storage Storage
	logger  logging.Logger
}

type Service interface {
	Create(ctx context.Context, recipe CreateRecipeDTO) (string, error)
	Patch(ctx context.Context, recipe EditRecipeDTO) error
	Get(ctx context.Context, id string) (Recipe, error)
	Delete(ctx context.Context, id string) error
}

func (s service) Delete(ctx context.Context, id string) error {
	err := s.storage.Delete(ctx, id)
	if err != nil {
		if _, ok := err.(*apperror.AppError); !ok {
			return apperror.NewAppError("Something went wrong in DB", apperror.WTC000001, err.Error())
		}
	}

	return err
}

func (s service) Get(ctx context.Context, id string) (Recipe, error) {
	r, err := s.storage.FindOne(ctx, id)
	if err != nil {
		if _, ok := err.(*apperror.AppError); !ok {
			return Recipe{}, apperror.NewAppError("Something went wrong in DB", apperror.WTC000001, err.Error())
		}
	}
	return r, err
}

func (s service) Patch(ctx context.Context, rDto EditRecipeDTO) error {
	var r Recipe

	// fill recipe model
	r = rDto.ToRecipe()

	err := s.storage.Update(ctx, r)
	if err != nil {
		if _, ok := err.(*apperror.AppError); !ok {
			return apperror.NewAppError("Something went wrong in DB", apperror.WTC000001, err.Error())
		}
	}

	return err
}

func (s service) Create(ctx context.Context, rDto CreateRecipeDTO) (string, error) {
	var r Recipe
	// fill recipe model
	r = rDto.ToRecipe()

	id, err := s.storage.Create(ctx, r)
	if err != nil {
		if _, ok := err.(*apperror.AppError); err != nil && !ok {
			return "", apperror.NewAppError("Something went wrong in DB", apperror.WTC000001, err.Error())
		}
	}
	return id, nil
}

func NewService(recipeStorage Storage, logger logging.Logger) (Service, error) {
	return &service{
		storage: recipeStorage,
		logger:  logger,
	}, nil
}
