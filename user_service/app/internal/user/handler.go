package user

import (
	"encoding/json"
	"fmt"
	"github.com/WTC-SYSTEM/apperror"
	"github.com/WTC-SYSTEM/logging"
	"github.com/WTC-SYSTEM/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
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

func (h *Handler) loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Do stuff here
		h.Logger.Info(fmt.Sprintf("URL:%s METHOD:%s", r.URL, r.Method))
		// Call the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		next.ServeHTTP(w, r)
	})
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc(usersURL, apperror.Middleware(h.GetUserByEmailAndPassword)).
		Methods(http.MethodGet)
	router.HandleFunc(usersURL, apperror.Middleware(h.CreateUser)).
		Methods(http.MethodPost)
	router.HandleFunc(userURL, apperror.Middleware(h.GetUser)).
		Methods(http.MethodGet)
	router.HandleFunc(usersURL, apperror.Middleware(h.UpdateUser)).
		Methods(http.MethodPatch)
	router.Use(h.loggingMiddleware)
	router.Use(h.CORS)
}

func (h *Handler) CreateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("CreateUser")

	w.Header().Set("Content-Type", "application/json")

	var userDto CreateUserDTO
	defer r.Body.Close()
	err := json.NewDecoder(r.Body).Decode(&userDto)
	if err != nil {
		return apperror.BadRequestError("Error deserializing JSON")
	}

	err = h.Validator.Struct(userDto)

	if err != nil {
		return apperror.BadRequestError("Error validating JSON")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userDto.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	UUID, err := h.UserService.Create(r.Context(), *userDto.Hashed(hashedPassword))
	if err != nil {
		return err
	}

	w.Header().Set("Location", fmt.Sprintf("%s/%s", usersURL, UUID))
	w.WriteHeader(http.StatusCreated)

	return nil
}

func (h *Handler) GetUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GetUser")

	w.Header().Set("Content-Type", "application/json")

	vars := mux.Vars(r)
	uuid := vars["uuid"]
	user, err := h.UserService.GetOne(r.Context(), uuid)
	if err != nil {
		return apperror.BadRequestError("User doesn't exist")

	}
	b, err := utils.CreateResponse(user)

	if err != nil {
		return err
	}
	w.Write(b)
	w.WriteHeader(http.StatusOK)
	return nil
}

func (h *Handler) GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("GetUserByEmailAndPassword")
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

	b, err := utils.CreateResponse(user)

	if err != nil {
		return err
	}
	w.Write(b)
	w.WriteHeader(http.StatusOK)

	return nil
}

func (h *Handler) UpdateUser(w http.ResponseWriter, r *http.Request) error {
	h.Logger.Info("UpdateUser")
	w.Header().Set("Content-Type", "application/json")

	var updateUserDto UpdateUserDTO

	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&updateUserDto); err != nil {
		return apperror.BadRequestError("Decoding to dto failed")
	}

	err := h.UserService.Update(r.Context(), updateUserDto)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusAccepted)

	return nil
}
