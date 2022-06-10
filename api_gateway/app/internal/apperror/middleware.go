package apperror

import (
	"encoding/json"
	"errors"
	"net/http"
)

type appHandler func(http.ResponseWriter, *http.Request) error

func Middleware(h appHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var appErr *AppError
		encoder := json.NewEncoder(w)
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
			encoder.Encode(ErrNotKnown)
		}
	}
}
