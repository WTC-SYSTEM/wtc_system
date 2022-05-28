package recipe

import (
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/hawkkiller/wtc_system/recipe_service/internal/apperror"
	"github.com/hawkkiller/wtc_system/recipe_service/pkg/logging"
	"net/http"
	"strconv"
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
	//var recipe Recipe
	//if err := json.NewDecoder(r.Body).Decode(&recipe); err != nil {
	//	h.Logger.Error(err)
	//	return apperror.BadRequestError("invalid json")
	//}

	err := r.ParseMultipartForm(32 << 20)
	if err != nil {
		return err
	}
	_, f, err := r.FormFile("image")
	if err != nil {
		return err
	}
	w.Write([]byte(strconv.FormatInt(f.Size, 10)))
	//if err := h.Validator.Struct(recipe); err != nil {
	//	h.Logger.Error(err)
	//	w.WriteHeader(http.StatusBadRequest)
	//	return
	//}

	//if err := h.RecipeService.CreateRecipe(recipe); err != nil {
	//	h.Logger.Error(err)
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	w.WriteHeader(http.StatusCreated)
	return nil
}
