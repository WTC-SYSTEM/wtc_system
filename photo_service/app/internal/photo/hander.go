package photo

import (
	"bytes"
	"github.com/WTC-SYSTEM/wtc_system/libs/apperror"
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
	"github.com/WTC-SYSTEM/wtc_system/libs/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"io"
	"mime/multipart"
	"net/http"
	"runtime/debug"
	"sync"
)

const (
	photo = "/api/v1/photo"
)

type handler struct {
	Logger    logging.Logger
	Validator *validator.Validate
	Service   PhotoService
}

// new handler
func NewHandler(logger logging.Logger, service PhotoService, validator *validator.Validate) *handler {
	return &handler{
		Logger:    logger,
		Validator: validator,
		Service:   service,
	}
}

func (h *handler) Register(router *mux.Router) {
	router.HandleFunc(photo, apperror.Middleware(h.UploadPhoto)).Methods(http.MethodPost)
}

func (h *handler) UploadPhoto(w http.ResponseWriter, r *http.Request) error {
	defer func() error {
		if err := recover(); err != nil {
			h.Logger.Error(err, debug.Stack())
			return err.(error)
		}
		return nil
	}()

	var wg sync.WaitGroup

	// parse multipart form
	err := r.ParseMultipartForm(5 << 20)
	if err != nil {
		return err
	}
	files := r.MultipartForm.File["images"]

	var photos = make([]*PhotoDTO, len(files))

	for i, f := range files {
		wg.Add(1)
		go func(f *multipart.FileHeader, i int) {
			defer wg.Done()

			openedFile, err := f.Open()

			if err != nil {
				panic(err)
			}

			buf := bytes.NewBuffer(nil)
			if _, err := io.Copy(buf, openedFile); err != nil {
				panic(err)
			}
			if err != nil {
				panic(err)
			}

			photos[i] = &PhotoDTO{
				Filename: f.Filename,
				Size:     f.Size,
				Bytes:    buf.Bytes(),
			}
			err = openedFile.Close()
			if err != nil {
				panic(err)
			}
		}(f, i)

	}

	wg.Wait()

	filteredPhotos := utils.Where(photos, func(p *PhotoDTO) bool {
		ext := http.DetectContentType(p.Bytes)
		return ext == "image/jpeg" || ext == "image/png"
	})

	if len(filteredPhotos) == 0 {
		return apperror.BadRequestError("no valid photo")
	}

	// get folder string from query
	folder := r.URL.Query().Get("folder")

	if folder == "" {
		return apperror.BadRequestError("no folder")
	}

	dto := &UploadDTO{Photos: filteredPhotos, Folder: folder}

	urls, err := h.Service.Upload(r.Context(), dto)
	if err != nil {
		return err
	}
	res, err := utils.CreateResponse(urls)
	if err != nil {
		return err
	}
	if _, err := w.Write(res); err != nil {
		return err
	}
	return nil
}
