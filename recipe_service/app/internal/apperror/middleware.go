package apperror

import (
	"encoding/json"
	"errors"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		encoder := json.NewEncoder(w)
		var appErr *AppError
		err := h(w, r)
		if err != nil {
			if errors.As(err, &appErr) {
				if errors.Is(err, ErrNotFound) {
					w.WriteHeader(http.StatusNotFound)
					encoder.Encode(err)
					return
				}
				err := err.(*AppError)
				w.WriteHeader(http.StatusBadRequest)
				encoder.Encode(err)
				return
			}
			w.WriteHeader(418)
			encoder.Encode(err.Error())
		}
	}
}
