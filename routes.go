package main

import "net/http"

type Route struct {
	Name        string
	Method      string
	Pattern     string
	HandlerFunc http.HandlerFunc
}

var routes = []Route{

	Route{
		"Index",
		"GET",
		"/",
		Index,
	},

	Route{
		"GetTicker",
		"POST",
		"/api/ticker/{pair}",
		GetTicker,
	},
}
