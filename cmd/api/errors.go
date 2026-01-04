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

func (app *application) editConflictResponse(w http.ResponseWriter, r *http.Request) {
	message := "unable to update the record due to an edit conflict, please try again"
	app.errorResponse(w, r, http.StatusConflict, message)
}

func (app *application) noUserFound(w http.ResponseWriter, r *http.Request) {
	message := "no user found with this ID"

	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) noRestaurantFound(w http.ResponseWriter, r *http.Request) {
	message := "no restaraunt found with this ID"

	app.errorResponse(w, r, http.StatusNotFound, message)
}

func (app *application) invalidUserCredintails(w http.ResponseWriter, r *http.Request) {
	message := "invalid credintails for this user"

	app.errorResponse(w, r, http.StatusUnprocessableEntity, message)
}
func (app *application) rateLimitExceededResponse(w http.ResponseWriter, r *http.Request) {
	message := "rate limit exceeded"
	app.errorResponse(w, r, http.StatusTooManyRequests, message)
}

func (app *application) invalidAuthenticationTokenResponse(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("WWW-Authentiacate", "Bearer")

	message := "invalid or missing authentication token"
	app.errorResponse(w, r, http.StatusUnauthorized, message)
}

func (app *application) authorizationRequierd(w http.ResponseWriter, r *http.Request) {
	message := "you don't have the right permission to visit this page"

	app.errorResponse(w, r, http.StatusUnauthorized, message)
}
