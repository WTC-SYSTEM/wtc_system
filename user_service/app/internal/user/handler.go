package user

import (
	"github.com/gorilla/mux"
	"github.com/hawkkiller/wtc_system/user_service/internal/apperror"
	"github.com/hawkkiller/wtc_system/user_service/pkg/logging"
	"net/http"
)

type Handler struct {
	Logger      logging.Logger
	UserService Service
}

func (h *Handler) Register(router *mux.Router) {
	router.HandleFunc("/", apperror.Middleware(h.Test)).Methods("GET")

}

func (h *Handler) Test(w http.ResponseWriter, r *http.Request) error {
	h.UserService.Create(r.Context(), CreateUserDTO{})
	return apperror.BadRequestError("no data")
}
