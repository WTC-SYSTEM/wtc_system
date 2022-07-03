package recipe

import "context"

type Storage interface {
	Create(ctx context.Context, r Recipe) error
	FindOne(ctx context.Context, id string) (Recipe, error)
	Update(ctx context.Context, r Recipe) error
	Delete(ctx context.Context, id string) error
}
