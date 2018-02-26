package main

import (
	"net/http"
)

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		HandleIndex,
	},
	Route{
		"UserShow",
		"GET",
		"/api/user/{userId}/info",
		HandleUserInfo,
	},
	Route{
		"UserCreate",
		"POST",
		"/api/user/create",
		HandleUserCreate,
	},
	Route{
		"UserGetAttribute",
		"GET",
		"/api/user/{userId}/{attribute}",
		HandleUserGetAttribute,
	},
	Route{
		"UserSetAttribute",
		"PUT",
		"/api/user/{userId}/{attribute}",
		HandleUserSetAttribute,
	},
}
