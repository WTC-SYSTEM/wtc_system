package recipe

import (
	"github.com/go-playground/validator/v10"
	"github.com/hawkkiller/wtc_system/recipe_service/pkg/logging"
)

type Handler struct {
	Logger        logging.Logger
	RecipeService Service
	Validator     *validator.Validate
}
