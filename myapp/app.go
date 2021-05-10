package myapp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// define user struct
type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
}

// temporary repository
var userMap map[int]*User
var lastId int

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

func updateUserInfoHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	// get updateUser from body
	updateUser := new(User)
	err := json.NewDecoder(request.Body).Decode(updateUser)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err)
		return
	}

	// check updateUser (by id) existence
	id, _ := strconv.Atoi(vars["id"])
	_, ok := userMap[id]
	if !ok {
		writer.WriteHeader(http.StatusNoContent)
		fmt.Fprint(writer, "No User ID:", id)
		return
	}

	userMap[id] = updateUser

	writer.WriteHeader(http.StatusOK)
}

func deleteUserInfoHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err)
		return
	}

	_, ok := userMap[id]
	if !ok {
		writer.WriteHeader(http.StatusOK)
		fmt.Fprint(writer, "No User ID:", id)
		return
	}

	delete(userMap, id)
	writer.WriteHeader(http.StatusOK)
	fmt.Fprint(writer, "Deleted User ID:", id)
}

func createUserHandler(writer http.ResponseWriter, request *http.Request) {

	// read
	user := new(User)
	err := json.NewDecoder(request.Body).Decode(user)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err)
		return
	}

	// process
	lastId++
	user.ID = lastId
	user.CreatedAt = time.Now()
	userMap[user.ID] = user
	respBody, err := json.Marshal(user)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		fmt.Fprint(writer, err)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(http.StatusCreated)
	fmt.Fprint(writer, string(respBody))
}

func getUserInfoHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err)
	}

	user, ok := userMap[id]
	if !ok {
		writer.WriteHeader(http.StatusOK)
		fmt.Fprint(writer, "No User ID:", id)
		return
	}

	aUserResposne, _ := json.Marshal(user)

	writer.Header().Set("Content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	fmt.Fprint(writer, string(aUserResposne))
}

func usersHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Get UserInfo bu /users/{id}")
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hello World")
}
