package main_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	s "strings"
	"testing"

	"github.com/gorilla/mux"

	"."
)

var a *mux.Router
var testUsers []main.User

func TestMain(m *testing.M) {
	a = main.NewRouter()
	testUsers = getTestUsers("test_users.json")

	os.Setenv("REDIS_URL", "localhost:6380")

	os.Exit(m.Run())

}

func TestUserCreate(t *testing.T) {

	main.InitializeRedis()

	testString, err := json.Marshal(testUsers[0])

	if err != nil {
		panic(err)
	}

	req, _ := http.NewRequest("POST", "/api/user/create", bytes.NewBuffer(testString))
	response := executeRequest(req)

	assertHttpResponse(t, response, 200, fmt.Sprintf("Successfully created user %s", testUsers[0].Username))
}

func TestUserNotCreatedTwice(t *testing.T) {

	main.InitializeRedis()

	testString, err := json.Marshal(testUsers[0])

	if err != nil {
		panic(err)
	}

	req1, _ := http.NewRequest("POST", "/api/user/create", bytes.NewBuffer(testString))
	response1 := executeRequest(req1)

	assertHttpResponse(t, response1, 200, fmt.Sprintf("Successfully created user %s", testUsers[0].Username))

	req2, _ := http.NewRequest("POST", "/api/user/create", bytes.NewBuffer(testString))
	response2 := executeRequest(req2)

	assertHttpResponse(t, response2, 409, "The username was already taken")

}

func TestUserCreateMultiple(t *testing.T) {

	main.InitializeRedis()

	testString1, err := json.Marshal(testUsers[0])

	if err != nil {
		panic(err)
	}

	req1, _ := http.NewRequest("POST", "/api/user/create", bytes.NewBuffer(testString1))
	response1 := executeRequest(req1)

	assertHttpResponse(t, response1, 200, fmt.Sprintf("Successfully created user %s", testUsers[0].Username))

	testString2, err := json.Marshal(testUsers[1])

	if err != nil {
		panic(err)
	}

	req2, _ := http.NewRequest("POST", "/api/user/create", bytes.NewBuffer(testString2))
	response2 := executeRequest(req2)

	assertHttpResponse(t, response2, 200, fmt.Sprintf("Successfully created user %s", testUsers[1].Username))

}

func TestUserInfo(t *testing.T) {

	main.InitializeRedis()

	TestUserCreate(t)

	req, _ := http.NewRequest("GET", "/api/user/0/info", nil)
	response := executeRequest(req)

	assertHttpResponse(t, response, 200, "{\"username\":\"mary\",\"name\":\"mary\",\"birthdate\":\"\",\"email\":\"\"}")
}

func TestUserInfoNotFound(t *testing.T) {

	main.InitializeRedis()

	req, _ := http.NewRequest("GET", "/api/user/234590782/info", nil)
	response := executeRequest(req)

	assertHttpResponse(t, response, 404, "This user id (234590782) couldn't be found!")
}

func TestUserAttributeUsername(t *testing.T) {

	main.InitializeRedis()
	TestUserCreate(t)

	req, _ := http.NewRequest("GET", "/api/user/0/username", nil)
	response := executeRequest(req)

	assertHttpResponse(t, response, 200, "mary")
}

func TestUserAttributeUnknownUser(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/user/234590782/username", nil)
	response := executeRequest(req)

	assertHttpResponse(t, response, 404, "Couldn't find user 234590782")
}

func TestUserAttributeBadUserId(t *testing.T) {

	req, _ := http.NewRequest("GET", "/api/user/gopher/username", nil)
	response := executeRequest(req)

	assertHttpResponse(t, response, 409, "Invalid userid gopher, please use integer")
}

func TestUserAttributeSetUsername(t *testing.T) {

	main.InitializeRedis()
	TestUserCreate(t)

	req1, _ := http.NewRequest("PUT", "/api/user/0/username", s.NewReader("marie"))
	response1 := executeRequest(req1)

	assertHttpResponse(t, response1, 200, "ok")

	req2, _ := http.NewRequest("GET", "/api/user/0/username", nil)
	response2 := executeRequest(req2)

	assertHttpResponse(t, response2, 200, "marie")
}

func TestUserAttributeSetUnknownUser(t *testing.T) {

	req, _ := http.NewRequest("PUT", "/api/user/234590782/username", s.NewReader("marie"))
	res := executeRequest(req)

	assertHttpResponse(t, res, 404, "Couldn't find user 234590782")
}

func TestUserAttributeSetInvalidUserID(t *testing.T) {

	req, _ := http.NewRequest("PUT", "/api/user/gopher/username", s.NewReader("marie"))
	res := executeRequest(req)

	assertHttpResponse(t, res, 409, "Invalid userid gopher, please use integer")
}

func TestUserAttributeSetEmptyAttribute(t *testing.T) {

	main.InitializeRedis()
	TestUserCreate(t)

	req, _ := http.NewRequest("PUT", "/api/user/0/username", s.NewReader(""))
	res := executeRequest(req)

	assertHttpResponse(t, res, 409, "Must pass in some body to /api/user/{userId}/{attribute}")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()

	a.ServeHTTP(rr, req)

	return rr
}

// assertHttpResponse wraps checks around the status code and response body
func assertHttpResponse(t *testing.T, r *httptest.ResponseRecorder, statusCode int, expectedBody string) {

	assertEqual(t, "response code", r.Result().StatusCode, statusCode)
	assertEqual(t, "response body", s.TrimSpace(r.Body.String()), expectedBody)

}

func assertEqual(t *testing.T, testType string, actual interface{}, expected interface{}) {

	if actual != expected {
		t.Errorf("The expected value for %s was %s, but received %s", testType, expected, actual)
	}
}

func getTestUsers(testFilename string) []main.User {

	f, err := ioutil.ReadFile(testFilename)

	if err != nil {
		panic(err)
	}

	var unpack []main.User
	err = json.Unmarshal(f, &unpack)
	if err != nil {
		panic(err)
	}

	return unpack
}
