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

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		app.errorResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed for this reousrce")

	})

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheck)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.userInformationHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.updateUserHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.deleteUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/restaurant", app.restaurantCreateHandler)

	return router
}
