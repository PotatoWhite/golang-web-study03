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

func TestUsers_withoutUsers(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	resp, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)
	respBody, _ := ioutil.ReadAll(resp.Body)
	temp := string(respBody)
	log.Println(temp)
	assert.Contains(string(respBody), "No Users")
}

func TestUsers_withUsers(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// create user #1
	user01 := new(User)
	user02 := new(User)

	user01.FirstName = "potato"
	user01.LastName = "white"
	user01.Email = "potato@gmail.com"

	user02.FirstName = "carrot"
	user02.LastName = "black"
	user02.Email = "carrot@example.com"

	user01dto, _ := json.Marshal(user01)
	user02dto, _ := json.Marshal(user01)

	resUser01, err1 := http.Post(ts.URL+"/users", "application/json", bytes.NewReader(user01dto))
	assert.NoError(err1)
	assert.Equal(http.StatusCreated, resUser01.StatusCode)

	resUser02, err2 := http.Post(ts.URL+"/users", "application/json", bytes.NewReader(user02dto))
	assert.NoError(err2)
	assert.Equal(http.StatusCreated, resUser02.StatusCode)

	respAll, err := http.Get(ts.URL + "/users")
	assert.NoError(err)
	assert.Equal(http.StatusOK, respAll.StatusCode)

	users := []*User{}
	err = json.NewDecoder(respAll.Body).Decode(&users)
	assert.NoError(err)
	assert.Equal(2, len(users))
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

func TestUpdateUser_nonexist(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// update fail
	updateUser := new(UpdatedUser)

	updateUser.FirstName = "carrot"
	updateUser.UpdatedFirstName = true

	updateUser.Email = "carrot@example.com"
	updateUser.UpdatedEmail = true

	updateUserReq, err := json.Marshal(updateUser)
	assert.NoError(err)

	req, err := http.NewRequest(http.MethodPut, ts.URL+"/users/1", bytes.NewReader(updateUserReq))
	req.Header.Add("Content-type", "application/json")
	resp, _ := http.DefaultClient.Do(req)
	assert.Equal(http.StatusNoContent, resp.StatusCode)

}

func TestUpdateUser(t *testing.T) {
	assert := assert.New(t)

	ts := httptest.NewServer(NewHandler())
	defer ts.Close()

	// create a User
	resp, err := http.Post(ts.URL+"/users", "application/json", strings.NewReader(`{"first_name":"potato", "last_name":"white", "email":"bravopotato@gmail.com"}`))
	assert.NoError(err)
	assert.Equal(http.StatusCreated, resp.StatusCode)

	createdUser := new(User)
	err = json.NewDecoder(resp.Body).Decode(createdUser)
	assert.NoError(err)
	assert.NotNil(createdUser.ID)

	// change properties
	createdUser.FirstName = "carrot"
	createdUser.Email = "carrot@example.com"

	updateUserReq, err := json.Marshal(createdUser)
	if err != nil {
		return
	}
	// request update user
	request, err := http.NewRequest(http.MethodPut, ts.URL+"/users/"+strconv.Itoa(createdUser.ID), bytes.NewReader(updateUserReq))
	resp, err = http.DefaultClient.Do(request)
	assert.NoError(err)
	assert.Equal(http.StatusOK, resp.StatusCode)

	// retrieve user to check
	respRetrieve, error := http.Get(ts.URL + "/users/" + strconv.Itoa(createdUser.ID))
	assert.NoError(error)

	retrievedUser := new(User)
	json.NewDecoder(respRetrieve.Body).Decode(retrievedUser)
	raw, err := ioutil.ReadAll(respRetrieve.Body)
	assert.NoError(err)
	log.Println(string(raw))

	assert.Equal("carrot", retrievedUser.FirstName)
	assert.Equal("carrot@example.com", retrievedUser.Email)
	assert.Equal("white", retrievedUser.LastName)

}
