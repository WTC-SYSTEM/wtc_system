package recipe

import (
	"encoding/json"
	"github.com/WTC-SYSTEM/apperror"
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"net/http"
	"runtime/debug"
)

const (
	recipe     = "/api/v1/recipe"
	editRecipe = "/api/v1/recipe/{id}"
)

type Handler struct {
	Logger        logging.Logger
	RecipeService Service
	Validator     *validator.Validate
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc(recipe, apperror.Middleware(h.CreateRecipe)).
		Methods(http.MethodPost)
	router.HandleFunc(recipe, apperror.Middleware(h.GetRecipe)).
		Methods(http.MethodGet)
	router.HandleFunc(editRecipe, apperror.Middleware(h.EditRecipe)).
		Methods(http.MethodPatch)
}

// CreateRecipe create recipe
func (h *Handler) CreateRecipe(w http.ResponseWriter, r *http.Request) error {
	defer func() {
		if err := recover(); err != nil {
			h.Logger.Error(err, debug.Stack())
		}
	}()
	var dto CreateRecipeDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("Failed to decode passed data")
	}

	if id, err := h.RecipeService.Create(r.Context(), dto); err != nil {
		return err
	} else {
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(id))
	}
	return nil
}

func (h *Handler) EditRecipe(w http.ResponseWriter, r *http.Request) error {
	defer func() {
		if err := recover(); err != nil {
			h.Logger.Error(err, debug.Stack())
		}
	}()
	var dto EditRecipeDTO

	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("Failed to decode passed data")
	}

	id := mux.Vars(r)["id"]

	if id == "" {
		return apperror.BadRequestError("Recipe id is empty")
	}

	dto.ID = id

	if err := h.RecipeService.Patch(r.Context(), dto); err != nil {
		return err
	}
	w.WriteHeader(http.StatusAccepted)
	return nil
}

func (h *Handler) GetRecipe(w http.ResponseWriter, r *http.Request) error {

	//defer func() {
	//	if err := recover(); err != nil {
	//		h.Logger.Error(err, debug.Stack())
	//	}
	//}()

	// get recipe id from url
	id := r.URL.Query().Get("id")

	// if id is empty return error
	if id == "" {
		return apperror.BadRequestError("Recipe id is empty")
	}

	// get recipe from storage
	recipe, err := h.RecipeService.Get(r.Context(), id)
	if err != nil {
		return err
	}

	// encode recipe to json and write to response
	if err := json.NewEncoder(w).Encode(recipe); err != nil {
		return err
	}

	w.WriteHeader(http.StatusOK)
	return nil
}
