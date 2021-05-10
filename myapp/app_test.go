package myapp

import (
	"bytes"
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

func TestIndex(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(respBody), "Hello")
}

func TestUsers(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(respBody), "Get UserInfo")
}

func TestGetUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users/89")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	assert.Contains(string(respBody), "No User ID:89")

	resp, err = http.Get(ts.URL + "/users/56")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	respBody, _ = ioutil.ReadAll(resp.Body)
	assert.Contains(string(respBody), "No User ID:56")
}

func TestCreateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"first_name":"potato", "last_name":"white", "email":"bravopotato@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	createdUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(createdUser)
	assert.NoError(err)
	assert.NotNil(createdUser.ID)

	resp, err = http.Get(ts.URL + "/users/" + strconv.Itoa(createdUser.ID))
	assert.NoError(err)
	retrieveUser := new(User)
	_ = json.NewDecoder(resp.Body).Decode(retrieveUser)

	assert.Equal(createdUser, retrieveUser)
}

func TestDeleteUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	request, err := http.NewRequest(http.MethodDelete, ts.URL+"/users/1", nil)
	assert.NoError(err)
	resp, errDo := http.DefaultClient.Do(request)
	assert.NoError(errDo)
	assert.Equal(http.StatusOK, resp.StatusCode)

	respBody, err := ioutil.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Contains(string(respBody), "No User ID:1")

	// create a user
	resp, err = http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"first_name":"potato", "last_name":"white", "email":"bravopotato@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	createdUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(createdUser)
	assert.NoError(err)
	assert.NotNil(createdUser.ID)

	// delete user
	request, err = http.NewRequest(http.MethodDelete, ts.URL+"/users/"+strconv.Itoa(createdUser.ID), nil)
	resp, err = http.DefaultClient.Do(request)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	respBody, err = ioutil.ReadAll(resp.Body)
	assert.NoError(err)
	assert.Contains(string(respBody), "Deleted User")
}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// update fail
	newRequest, err := http.NewRequest(http.MethodPut, ts.URL+"/users/1", strings.NewReader(`{"first_name":"potato", "last_name":"white", "email":"bravopotato@gmail.com"}`))
	newRequest.Header.Add("Content-type", "application/json")
	resp, _ := http.DefaultClient.Do(newRequest)
	assert.Equal(http.StatusNoContent, resp.StatusCode)

	// create a User
	resp, err = http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"first_name":"potato", "last_name":"white", "email":"bravopotato@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)
	createdUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(createdUser)
	assert.NoError(err)
	assert.NotNil(createdUser.ID)

	// change properties
	createdUser.FirstName = "carrot"
	createdUser.Email = "carrot@example.com"

	updateUser, err := json.Marshal(createdUser)
	if err != nil {
		return
	}
	// request update user
	request, err := http.NewRequest(http.MethodPut, ts.URL+"/users/"+strconv.Itoa(createdUser.ID), bytes.NewReader(updateUser) )
	resp, err = http.DefaultClient.Do(request)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// retrieve user to check
	respRetrieve, error := http.Get(ts.URL+"/users/"+strconv.Itoa(createdUser.ID),)
	assert.NoError(error)

	retrievedUser := new(User)
	json.NewDecoder(respRetrieve.Body).Decode(retrievedUser)
	raw, err := ioutil.ReadAll(respRetrieve.Body)
	assert.NoError(err)
	log.Println(string(raw))

	assert.Equal(retrievedUser, createdUser)
}
