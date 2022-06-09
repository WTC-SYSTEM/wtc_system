package db

import (
	"context"
	"github.com/hawkkiller/wtc_system/recipe_service/internal/recipe"
	"github.com/hawkkiller/wtc_system/recipe_service/pkg/client/postgresql"
	"github.com/hawkkiller/wtc_system/recipe_service/pkg/logging"
)

type db struct {
	client postgresql.Client
	logger logging.Logger
}

func NewStorage(c postgresql.Client, l logging.Logger) recipe.Storage {
	return &db{
		client: c,
		logger: l,
	}
}

func (d db) Create(ctx context.Context, recipe recipe.Recipe) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d db) FindOne(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}

func (d db) Update(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (d db) Delete(ctx context.Context, id string) error {
	//TODO implement me
	panic("implement me")
}
