package recipe

import (
	"bytes"
	"encoding/json"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/internal/apperror"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/logging"
	"github.com/WTC-SYSTEM/wtc_system/recipe_service/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"runtime/debug"
	"strconv"
	"strings"
)

const (
	recipe = "/api/v1/recipe"
)

type Handler struct {
	Logger        logging.Logger
	RecipeService Service
	Validator     *validator.Validate
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc(recipe, apperror.Middleware(h.CreateRecipe)).
		Methods(http.MethodPost)
}

// CreateRecipe create recipe
func (h *Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) error {
	defer func() {
		if err := recover(); err != nil {
			h.Logger.Error(err, debug.Stack())
		}
	}()
	var dto CreateRecipeDTO

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}
	var stepsCount int
	for s := range r.MultipartForm.Value {
		if strings.Contains(s, "step") {
			stepsCount++
		}
	}

	// get recipe-data and decode json in multipart form
	err = func() error {
		v := r.MultipartForm.Value["recipe-data"][0]
		if err != nil {
			return err
		}
		d := json.NewDecoder(bytes.NewReader([]byte(v)))
		err = d.Decode(&dto)
		if err != nil {
			return err
		}
		return nil
	}()
	if err != nil {
		return err
	}

	for _, f := range r.MultipartForm.File["recipe-photos"] {
		err := func() error {
			open, err := f.Open()
			if err != nil {
				return err
			}
			defer func(open multipart.File) {
				err := open.Close()
				if err != nil {
					h.Logger.Errorf(err.Error())
				}
			}(open)
			r, err := ioutil.ReadAll(open)
			if err != nil {
				h.Logger.Errorf(err.Error())
				return err
			}

			dto.Photos = append(dto.Photos, Photo{
				MimeType: http.DetectContentType(r),
				Filename: utils.FileNameWithoutExtSliceNotation(f.Filename),
				Size:     f.Size,
				Data:     r,
			})
			return nil
		}()
		if err != nil {
			return err
		}
	}
	for i := 0; i < stepsCount; i++ {
		dto.Steps = append(dto.Steps, CreateStepDTO{})
		// read step-i-data and decode json in multipart form
		err := func() error {
			v := r.MultipartForm.Value["step-"+strconv.Itoa(i+1)+"-data"][0]
			if err != nil {
				return err
			}
			d := json.NewDecoder(bytes.NewReader([]byte(v)))
			err = d.Decode(&dto.Steps[i])
			if err != nil {
				return err
			}
			return nil
		}()
		if err != nil {
			return err
		}
		for _, f := range r.MultipartForm.File["step-"+strconv.Itoa(i+1)+"-photos"] {
			err := func() error {
				open, err := f.Open()
				if err != nil {
					return err
				}
				defer func(open multipart.File) {
					err := open.Close()
					if err != nil {
						h.Logger.Errorf(err.Error())
					}
				}(open)
				r, err := ioutil.ReadAll(open)
				if err != nil {
					return err
				}
				dto.Steps[i].Photos = append(dto.Steps[i].Photos, Photo{
					MimeType: http.DetectContentType(r),
					Filename: utils.FileNameWithoutExtSliceNotation(f.Filename),
					Size:     f.Size,
					Data:     r,
				})
				return nil
			}()
			if err != nil {
				return err
			}
		}

	}
	if err := h.RecipeService.Create(r.Context(), dto); err != nil {
		return err
	}
	w.WriteHeader(http.StatusCreated)
	return nil
}
