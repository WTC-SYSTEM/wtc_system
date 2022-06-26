package recipe

import (
	"context"
	"encoding/base64"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/logging"
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

	// upload recipe photos
	for _, p := range rDto.Photos {
		b, err := base64.StdEncoding.DecodeString(p)
		if err != nil {
			return err
		}
		url, err := s.storage.UploadFile(ctx, b, "recipe-photos")
		if err != nil {
			return err
		}
		r.Photos = append(r.Photos, url)
	}

	// upload steps photos and fill steps
	for i, sDto := range rDto.Steps {
		r.Steps = append(r.Steps, sDto.ToStep())
		for _, p := range sDto.Photos {
			b, err := base64.StdEncoding.DecodeString(p)
			if err != nil {
				return err
			}
			url, err := s.storage.UploadFile(ctx, b, "recipe-steps-photos")
			if err != nil {
				return err
			}
			r.Steps[i].Photos = append(r.Steps[i].Photos, url)
		}
	}

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

	// upload recipe photos
	for _, p := range rDto.Photos {
		b, err := base64.StdEncoding.DecodeString(p)
		if err != nil {
			return err
		}
		url, err := s.storage.UploadFile(ctx, b, "recipe-photos")
		if err != nil {
			return err
		}
		r.Photos = append(r.Photos, url)
	}

	// upload steps photos and fill steps
	for i, sDto := range rDto.Steps {
		r.Steps = append(r.Steps, sDto.ToStep())
		for _, p := range sDto.Photos {
			b, err := base64.StdEncoding.DecodeString(p)
			if err != nil {
				return err
			}
			url, err := s.storage.UploadFile(ctx, b, "recipe-steps-photos")
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
