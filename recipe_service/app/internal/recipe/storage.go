package recipe

import "context"

type Storage interface {
	Create(ctx context.Context, r Recipe) error
	UploadFile(ctx context.Context, photo Photo, folder string) (string, error)
	FindOne(ctx context.Context, id string) error
	Update(ctx context.Context) error
	Delete(ctx context.Context, id string) error
}
