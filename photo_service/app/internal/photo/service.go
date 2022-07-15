package photo

import (
	"context"
	"github.com/WTC-SYSTEM/logging"
	"sync"
)

type service struct {
	storage Storage
	logger  logging.Logger
}

type PhotoService interface {
	Upload(ctx context.Context, dto *UploadDTO) ([]UploadedItem, error)
}

func NewService(storage Storage, logger logging.Logger) PhotoService {
	return &service{
		storage: storage,
		logger:  logger,
	}
}

func (s service) Upload(ctx context.Context, dto *UploadDTO) ([]UploadedItem, error) {
	var wg sync.WaitGroup

	var photos = make([]UploadedItem, len(dto.Photos))

	for i, photo := range dto.Photos {
		wg.Add(1)
		go func(i int, photo *PhotoDTO) {
			defer wg.Done()
			url, err := s.storage.Create(ctx, photo.Bytes, dto.Folder)
			if err != nil {
				panic(err)
			}
			photos[i].Url = url
			photos[i].Filename = photo.Filename
		}(i, photo)
	}
	wg.Wait()
	return photos, nil
}
