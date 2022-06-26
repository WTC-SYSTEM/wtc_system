package recipe

import "context"

type Storage interface {
	Create(ctx context.Context, r Recipe) error
	UploadFile(ctx context.Context, photo []byte, folder string) (string, error)
	FindOne(ctx context.Context, id string) (Recipe, error)
	Update(ctx context.Context, r Recipe) error
	Delete(ctx context.Context, id string) error
}
