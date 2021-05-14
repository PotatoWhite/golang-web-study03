package myapp

import (
	"github.com/gorilla/mux"
	"net/http"
)

// NewHandler make a myapp Handler
func NewHandler() http.Handler {
	userMap = make(map[int]*User)
	lastId = 0

	mux := mux.NewRouter()
	mux.HandleFunc("/", indexHandler)
	mux.HandleFunc("/users", usersHandler).Methods(http.MethodGet)
	mux.HandleFunc("/users", createUserHandler).Methods(http.MethodPost)
	mux.HandleFunc("/users/{id:[0-9a-z]+}", getUserInfoHandler).Methods(http.MethodGet)
	mux.HandleFunc("/users/{id:[0-9a-z]+}", deleteUserInfoHandler).Methods(http.MethodDelete)
	mux.HandleFunc("/users/{id:[0-9a-z]+}", updateUserInfoHandler).Methods(http.MethodPut)
	return mux
}
