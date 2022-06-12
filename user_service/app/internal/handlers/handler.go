package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handler interface {
	Register(router *mux.Router)
	CreateUser(w http.ResponseWriter, r *http.Request) error
	GetUser(w http.ResponseWriter, r *http.Request) error
	GetUserByEmailAndPassword(w http.ResponseWriter, r *http.Request) error
	UpdateUser(w http.ResponseWriter, r *http.Request) error
}
