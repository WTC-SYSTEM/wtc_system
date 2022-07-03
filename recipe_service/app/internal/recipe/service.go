package recipe

import (
	"context"
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
)

type service struct {
	storage Storage
	logger  logging.Logger
}

type Service interface {
	Create(ctx context.Context, recipe CreateRecipeDTO) error
	Patch(ctx context.Context, recipe EditRecipeDTO) error
	Get(ctx context.Context, id string) (Recipe, error)
}

func (s service) Get(ctx context.Context, id string) (Recipe, error) {
	return s.storage.FindOne(ctx, id)
}

func (s service) Patch(ctx context.Context, rDto EditRecipeDTO) error {
	var r Recipe

	// fill recipe model
	r = rDto.ToRecipe()

	err := s.storage.Update(ctx, r)
	if err != nil {
		return err
	}
	return nil
}

func (s service) Create(ctx context.Context, rDto CreateRecipeDTO) error {
	var r Recipe
	// fill recipe model
	r = rDto.ToRecipe()

	err := s.storage.Create(ctx, r)
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
