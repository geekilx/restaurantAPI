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
		Role      string `json:"role"`
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

	err = app.models.Users.Insert(&user)
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

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetUser(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
		Email     string `json:"email"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if input.FirstName != "" {
		user.FirstName = input.FirstName
	}
	if input.LastName != "" {
		user.LastName = input.LastName
	}
	if input.Email != "" {
		user.Email = input.Email
	}

	//TODO: adding a forget password functionalitty
	// if input.Password != "" {
	// 	user.Password.Set(input.Password)
	// }

	err = app.models.Users.Update(*user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrConflictEdit):
			app.editConflictResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"message": "user successfully updated"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) deleteUserHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	//TODO: adding passowrd check after adding user authentication
	err = app.models.Users.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"messeage": "user successfully deleted."}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) userInformationHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(w, r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetUser(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.noUserFound(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	//TODO: adding user authentication in order to keep other users from seeing another user informations is URGENT
	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"user": user}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
