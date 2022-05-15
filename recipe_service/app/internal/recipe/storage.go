package recipe

import "context"

type Storage interface {
	Create(ctx context.Context) (string, error)
	FindOne(ctx context.Context, id string) error
	Update(ctx context.Context) error
	Delete(ctx context.Context, id string) error
}
