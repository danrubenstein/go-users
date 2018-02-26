package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
)

func HandleIndex(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome!\n")
}

func HandleUserInfo(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	userId, err := strconv.Atoi(v["userId"])

	if err != nil {
		http.Error(w, err.Error(), 409)
	}

	user, err := ClientUserInfo(userId)
	if err != nil {

		if strings.Contains(err.Error(), "couldn't be found!") {
			http.Error(w, err.Error(), 404)
		} else {
			http.Error(w, err.Error(), 500)
		}
	} else {

		responseBytes, err := json.Marshal(user)
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
		fmt.Fprintln(w, string(responseBytes))
	}
}

func HandleUserCreate(w http.ResponseWriter, r *http.Request) {

	var user User

	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&user)

	if err != nil {
		fmt.Printf("error - %s", err)
		http.Error(w, "JSON invalid", 409)
		return
	}

	if err = ClientUserCreate(user); err != nil {
		http.Error(w, err.Error(), 409)
	} else {
		fmt.Fprintln(w, fmt.Sprintf("Successfully created user %s", user.Username))
	}
}

func HandleUserGetAttribute(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	userId, err := strconv.Atoi(v["userId"])

	if err != nil {
		errorMessage := fmt.Sprintf("Invalid userid %s, please use integer", v["userId"])
		http.Error(w, errorMessage, 409)
		return
	}

	attr, err := ClientUserGetAttribute(userId, v["attribute"])

	if err != nil {
		if strings.Contains(err.Error(), "Couldn't find") {
			http.Error(w, err.Error(), 404)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	fmt.Fprintln(w, attr)
}

func HandleUserSetAttribute(w http.ResponseWriter, r *http.Request) {

	v := mux.Vars(r)
	userId, err := strconv.Atoi(v["userId"])

	if err != nil {
		errorMessage := fmt.Sprintf("Invalid userid %s, please use integer", v["userId"])
		http.Error(w, errorMessage, 409)
		return
	}

	requestBody, err := ioutil.ReadAll(r.Body)

	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	if len(requestBody) == 0 {
		http.Error(w, "Must pass in some body to /api/user/{userId}/{attribute}", 409)
		return
	}

	err = ClientUserSetAttribute(userId, v["attribute"], string(requestBody))

	if err != nil {
		if strings.Contains(err.Error(), "Couldn't find") {
			http.Error(w, err.Error(), 404)
		} else {
			http.Error(w, err.Error(), 500)
		}
		return
	}

	fmt.Fprintln(w, "ok")
}
