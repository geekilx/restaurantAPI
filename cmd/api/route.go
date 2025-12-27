package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) route() http.Handler {
	router := httprouter.New()

	// changing the default not found response in httprouter
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFoundResponse(w, r)
	})

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheck)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)

	return router
}
