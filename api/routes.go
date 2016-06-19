package main

import "net/http"

// Route is a http API route
type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

// Routes is an array of Route objects
type Routes []Route

var routes = Routes{
	Route{
		"Index",
		"GET",
		"/",
		Index,
	},
	Route{
		"VBoxVersion",
		"GET",
		"/virtualbox/version",
		VBoxVersion,
	},
	Route{
		"VBoxList",
		"GET",
		"/virtualbox/list",
		VBoxList,
	},
}
