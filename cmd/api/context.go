package main

import (
	"context"
	"net/http"

	"github.com/geekilx/restaurantAPI/internal/models"
)

type contextKey string

const userContextKey = contextKey("user")

func (app *application) getUserContext(r *http.Request) *models.User {

	user, ok := r.Context().Value(userContextKey).(*models.User)
	if !ok {
		panic("missing user value in the request context")
	}

	return user

}

func (app *application) setUserContext(w http.ResponseWriter, r *http.Request, user *models.User) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)

	return r.WithContext(ctx)

}
