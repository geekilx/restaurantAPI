package main

import (
	"errors"
	"net/http"

	"github.com/geekilx/restaurantAPI/internal/models"
	"github.com/geekilx/restaurantAPI/internal/validator"
)

func (app *application) createUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
		Password  string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var user = models.User{
		FirstName: input.FirstName,
		LastName:  input.LastName,
		Email:     input.Email,
	}

	v := validator.New()

	if models.ValidateUsers(v, user, input.Password); !v.Valid() {
		app.failedValidationResponse(w, r, v)
		return
	}

	err = user.Password.Set(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.UserModel.Insert(&user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			app.errorResponse(w, r, http.StatusUnprocessableEntity, "duplicate email is not premitted")
		default:
			app.serverErrorResponse(w, r, err)

		}
		return
	}

	err = app.writeJSON(w, r, http.StatusCreated, jsFmt{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
