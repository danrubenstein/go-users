package main_test

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/http/httptest"
    "os"
    "testing"

    "github.com/gorilla/mux"
    
    
    "."
)

var a *mux.Router
var testUsers []main.User

func TestMain(m *testing.M) {
    a = main.NewRouter()
    testUsers = getTestUsers("test_users.json")
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

    assertEqual(t, "response code", response.Result().StatusCode, 200)
    assertEqual(t, "response body", response.Body.String(), fmt.Sprintf("Successfully created user %s\n", testUsers[0].Username))
}

func TestUserNotCreatedTwice(t *testing.T) { 

    main.InitializeRedis()

    testString, err := json.Marshal(testUsers[0])
    
    if err != nil { 
        panic(err)
    }

    req1, _ := http.NewRequest("POST", "/api/user/create", bytes.NewBuffer(testString))
    response1 := executeRequest(req1)

    assertEqual(t, "response code", response1.Result().StatusCode, 200)
    assertEqual(t, "response body", response1.Body.String(), fmt.Sprintf("Successfully created user %s\n", testUsers[0].Username))

    req2, _ := http.NewRequest("POST", "/api/user/create", bytes.NewBuffer(testString))
    response2 := executeRequest2(req2)

    assertEqual(t, "response code", response2.Result().StatusCode, 409)
    assertEqual(t, "response body", response2.Body.String(), "The username was already taken\n")

}

func TestUserCreateMultiple(t *testing.T) { 

    main.InitializeRedis()

    testString1, err := json.Marshal(testUsers[0])
    
    if err != nil { 
        panic(err)
    }

    req1, _ := http.NewRequest("POST", "/api/user/create", bytes.NewBuffer(testString1))
    response1 := executeRequest(req1)

    assertEqual(t, "response code", response1.Result().StatusCode, 200)
    assertEqual(t, "response body", response1.Body.String(), fmt.Sprintf("Successfully created user %s\n", testUsers[0].Username))

    testString2, err := json.Marshal(testUsers[1])
    
    if err != nil { 
        panic(err)
    }


    req2, _ := http.NewRequest("POST", "/api/user/create", bytes.NewBuffer(testString2))
    response2 := executeRequest2(req2)

    assertEqual(t, "response code", response2.Result().StatusCode, 200)
    assertEqual(t, "response body", response2.Body.String(), fmt.Sprintf("Successfully created user %s\n", testUsers[1].Username))

}


func executeRequest(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()

    a.ServeHTTP(rr, req)

    return rr
}

func executeRequest2(req *http.Request) *httptest.ResponseRecorder {
    rr := httptest.NewRecorder()

    a.ServeHTTP(rr, req)

    return rr
}

func assertEqual(t *testing.T, testType string, actual interface{}, expected interface{}) {
    if actual != expected {
        t.Errorf("The expected value for %s was %s, but received %s", testType, expected, actual)
    }
}

func getTestUsers(testFilename string) []main.User{ 

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