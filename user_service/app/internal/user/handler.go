package user

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/hawkkiller/wtc_system/user_service/internal/apperror"
	"github.com/hawkkiller/wtc_system/user_service/pkg/logging"
	"github.com/hawkkiller/wtc_system/user_service/pkg/utils"
	"net/http"
)

const (
	usersURL = "/api/users"
	userURL  = "/api/users/:uuid"
)

type Handler struct {
	Logger      logging.Logger
	UserService Service
	Validator   *validator.Validate
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc(usersURL, apperror.Middleware(h.CreateUser)).Methods("POST")
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var userDto CreateUserDTO

	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		return apperror.BadRequestError("Error deserializing JSON")
	}

	defer r.Body.Close()
	err = h.Validator.Struct(userDto)
	if err != nil {
		return apperror.BadRequestError("Error validating JSON")
	}

	err = h.UserService.Create(r.Context(), userDto)
	if err != nil {
		return err
	}

	b, err := utils.CreateResponse(map[string]any{
		"message": "Successfully registered user",
		"code":    RegSuccess,
	})

	if err != nil {
		return err
	}
	w.Write(b)
	w.WriteHeader(200)
	return nil
}
