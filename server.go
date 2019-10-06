package main

import (
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {
	RunServer()
}

func RunServer() {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range routes {
		var handler http.Handler

		handler = route.HandlerFunc
		handler = Logger(handler, route.Name)

		router.
			Methods(route.Method).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)
	}

	log.Fatal(http.ListenAndServe(":8080", router))
}
