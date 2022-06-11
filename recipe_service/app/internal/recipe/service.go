package recipe

import (
	"context"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/logging"
)

type service struct {
	storage Storage
	logger  logging.Logger
}

type Service interface {
	Create(ctx context.Context, recipe CreateRecipeDTO) error
}

func (s service) Create(ctx context.Context, rDto CreateRecipeDTO) error {
	var r Recipe

	// fill recipe model
	r = rDto.ToRecipe()

	// upload recipe photos
	for _, p := range rDto.Photos {
		url, err := s.storage.UploadFile(ctx, p, "recipe-photos")
		if err != nil {
			return err
		}
		r.Photos = append(r.Photos, url)
	}

	// upload steps photos and fill steps
	for i, sDto := range rDto.Steps {
		r.Steps = append(r.Steps, sDto.ToStep())
		for _, p := range sDto.Photos {
			url, err := s.storage.UploadFile(ctx, p, "recipe-steps-photos")
			if err != nil {
				return err
			}
			r.Steps[i].Photos = append(r.Steps[i].Photos, url)
		}
	}

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
