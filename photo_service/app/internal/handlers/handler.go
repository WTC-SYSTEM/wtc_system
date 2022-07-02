package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handler interface {
	Register(router *mux.Router)
	UploadPhoto(w http.ResponseWriter, r *http.Request) error
}
