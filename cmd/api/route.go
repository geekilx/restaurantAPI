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
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.updateUserHandler)

	return router
}
