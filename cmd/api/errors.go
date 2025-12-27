package main

import (
	"net/http"

	"github.com/geekilx/restaurantAPI/internal/validator"
)

func (app *application) logError(r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)

}

func (app *application) errorResponse(w http.ResponseWriter, r *http.Request, status int, message any) {
	msg := jsFmt{"error": message}
	err := app.writeJSON(w, r, status, msg, nil)
	if err != nil {
		app.logError(r, err)
		w.WriteHeader(500)
	}

}

func (app *application) serverErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logError(r, err)

	msg := "server encountered an erorr and could not process your request"

	app.errorResponse(w, r, http.StatusInternalServerError, msg)
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request) {
	msg := "the requested resource could not be found"

	app.errorResponse(w, r, http.StatusNotFound, msg)
}

func (app *application) failedValidationResponse(w http.ResponseWriter, r *http.Request, v *validator.Validator) {
	app.errorResponse(w, r, http.StatusUnprocessableEntity, v.FieldErorrs)
}
