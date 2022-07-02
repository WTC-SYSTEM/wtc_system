package auth

import (
	"encoding/json"
	"fmt"
	"github.com/WTC-SYSTEM/wtc_system/api_gateway/internal/apperror"
	"github.com/WTC-SYSTEM/wtc_system/api_gateway/internal/client/user_service"
	"github.com/WTC-SYSTEM/wtc_system/api_gateway/pkg/jwt"
	"github.com/WTC-SYSTEM/wtc_system/libs/logging"
	"github.com/gorilla/mux"
	"net/http"
)

const (
	authURL   = "/api/auth"
	signupURL = "/api/signup"
)

type Handler struct {
	Logger      logging.Logger
	UserService user_service.UserService
	JWTHelper   jwt.Helper
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
	router.HandleFunc(authURL, apperror.Middleware(h.Auth)).Methods(http.MethodPost, http.MethodPut)
	router.HandleFunc(signupURL, apperror.Middleware(h.Signup)).Methods(http.MethodPost)
	fs := http.FileServer(http.Dir("../../internal/public"))
	router.PathPrefix("/public/").Handler(http.StripPrefix("/public/", fs)).
		Methods("GET")
	router.Use(h.loggingMiddleware)
	router.Use(h.CORS)
}

func (h *Handler) Signup(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	defer r.Body.Close()
	var dto user_service.CreateUserDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		return apperror.BadRequestError("failed to decode data")
	}

	u, err := h.UserService.Create(r.Context(), dto)
	if err != nil {
		return err
	}
	token, err := h.JWTHelper.GenerateAccessToken(u)
	if err != nil {
		return err
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(token)

	return nil
}

func (h *Handler) Auth(w http.ResponseWriter, r *http.Request) error {
	w.Header().Set("Content-Type", "application/json")

	var token []byte
	var err error
	switch r.Method {
	case http.MethodPost:
		defer r.Body.Close()
		var dto user_service.SignInUserDTO
		if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
			return apperror.BadRequestError("failed to decode data")
		}
		u, err := h.UserService.GetByEmailAndPassword(r.Context(), dto.Email, dto.Password)
		if err != nil {
			return err
		}
		token, err = h.JWTHelper.GenerateAccessToken(u)
		if err != nil {
			return err
		}
	case http.MethodPut:
		defer r.Body.Close()
		var rt jwt.RT
		if err := json.NewDecoder(r.Body).Decode(&rt); err != nil {
			return apperror.BadRequestError("failed to decode data")
		}
		token, err = h.JWTHelper.UpdateRefreshToken(rt)
		if err != nil {
			return err
		}
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(token)

	return err
}
