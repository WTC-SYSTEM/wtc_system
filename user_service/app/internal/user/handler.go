package user

import (
	"encoding/json"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/hawkkiller/wtc_system/user_service/internal/apperror"
	"github.com/hawkkiller/wtc_system/user_service/pkg/logging"
	"github.com/hawkkiller/wtc_system/user_service/pkg/utils"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

const (
	usersURL = "/api/v1/users"
	userURL  = "/api/v1/users/{uuid}"
)

type Handler struct {
	Logger      logging.Logger
	UserService Service
	Validator   *validator.Validate
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc(usersURL, apperror.Middleware(h.CreateUser)).
		Methods("POST")
	router.HandleFunc(usersURL, apperror.Middleware(h.GetUserByEmailAndPassword)).
		Methods("GET")
	router.HandleFunc(userURL, apperror.Middleware(h.GetUser)).
		Methods("GET")
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	var userDto CreateUserDTO
	w.Header().Set("Content-Type", "application/json")

	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		return apperror.BadRequestError("Error deserializing JSON")
	}
	defer r.Body.Close()

	err = h.Validator.Struct(userDto)

	if err != nil {
		return apperror.BadRequestError("Error validating JSON")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	err = h.UserService.Create(r.Context(), *userDto.Hashed(hashedPassword))
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

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	uuid := vars["uuid"]
	user, err := h.UserService.GetOne(r.Context(), uuid)
	if err != nil {
		return apperror.BadRequestError("User doesn't exist")

	}
	b, err := utils.CreateResponse(map[string]any{
		"message": "Successfully got user",
		"code":    GetUserSuccess,
		"user":    user,
	})

	if err != nil {
		return err
	}
	w.Write(b)
	w.WriteHeader(200)
	return nil
}

func (h *Handler) GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	email := r.URL.Query().Get("email")
	password := r.URL.Query().Get("password")
	if email == "" || password == "" {
		return apperror.BadRequestError("Email or password is empty")
	}
	dto := GetUserByEmailAndPasswordDTO{Password: password, Email: email}
	user, err := h.UserService.GetByEmailAndPassword(r.Context(), dto)

	if err != nil {
		if _, ok := err.(*apperror.AppError); ok {
			return err
		}
		return apperror.BadRequestError("User doesn't exist")
	}

	b, err := utils.CreateResponse(map[string]any{
		"message": "Successfully got user",
		"code":    GetUserSuccess,
		"user":    user,
	})

	if err != nil {
		return err
	}
	w.Write(b)
	w.WriteHeader(200)

	return nil
}
