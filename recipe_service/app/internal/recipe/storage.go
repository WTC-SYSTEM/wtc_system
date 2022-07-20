package recipe

import "context"

type Storage interface {
	Create(ctx context.Context, r Recipe) (string, error)
	FindOne(ctx context.Context, id string) (Recipe, error)
	FindAll(ctx context.Context, filter Filter) ([]Recipe, error)
	Update(ctx context.Context, r Recipe) error
	Delete(ctx context.Context, id string) error
}
