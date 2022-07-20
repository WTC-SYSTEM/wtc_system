package handlers

import (
	"github.com/gorilla/mux"
	"net/http"
)

type Handler interface {
	Register(router *mux.Router)
	CreateRecipe(w http.ResponseWriter, r *http.Request) error
	EditRecipe(w http.ResponseWriter, r *http.Request) error
	GetRecipe(w http.ResponseWriter, r *http.Request) error
	DeleteRecipe(w http.ResponseWriter, r *http.Request) error
	GetRecipes(w http.ResponseWriter, r *http.Request) error
}
