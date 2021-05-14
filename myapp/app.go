package myapp

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"time"
)

// temporary repository
var userMap map[int]*User
var lastId int

func updateUserInfoHandler(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)

	// get requestUser from body
	requestUser := new(User)
	err := json.NewDecoder(request.Body).Decode(requestUser)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(writer, err)
		return
	}

	// check requestUser (by id) existence
	id, _ := strconv.Atoi(vars["id"])
	targetUser, ok := userMap[id]
	if !ok {
		writer.WriteHeader(http.StatusNoContent)
		fmt.Fprint(writer, "No User ID:", id)
		return
	}

	if requestUser.FirstName != "" {
		targetUser.FirstName = requestUser.FirstName
	}

	if requestUser.LastName != "" {
		targetUser.LastName = requestUser.LastName
	}

	if requestUser.Email != "" {
		targetUser.Email = requestUser.Email
	}

	userMap[id] = targetUser

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
	if len(userMap) == 0 {
		writer.WriteHeader(http.StatusOK)
		fmt.Fprint(writer, "No Users")
	}

	users := []*User{}
	for _, element := range userMap {
		users = append(users, element)
	}

	response, _ := json.Marshal(users)
	writer.Header().Set("Content-type", "application/json")
	writer.WriteHeader(http.StatusOK)
	fmt.Fprint(writer, string(response))
}

func indexHandler(writer http.ResponseWriter, request *http.Request) {
	fmt.Fprint(writer, "Hello World")
}
