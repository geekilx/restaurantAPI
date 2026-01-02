package main

import (
	"errors"
	"net/http"
	"time"

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

	token, err := app.models.Tokens.New(72*time.Hour, user.ID, models.ActivationScope)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	data := map[string]any{
		"userID":          user.ID,
		"activationToken": token.PlainHash,
	}

	err = app.mailer.Send(user.Email, "template.tmpl", data)
	if err != nil {
		app.logger.Error(err.Error())
		return
	}

	err = app.models.Permissions.AddForUser(user.ID, "restaurant:read")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, r, http.StatusCreated, jsFmt{"user": user, "message": "Please check your email in order to activate your account"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) updateUserHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
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

	err = app.models.Users.Update(user)
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
	id, err := app.readIDParam(r)
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

	id, err := app.readIDParam(r)
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

func (app *application) userActivateHandler(w http.ResponseWriter, r *http.Request) {

	var input struct {
		Token string `json:"token"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	id, err := app.models.Tokens.GetByToken(input.Token)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			v.AddError("token", "invalid or expired activation token")
			app.failedValidationResponse(w, r, v)
		default:
			app.serverErrorResponse(w, r, err)
		}
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

	user.IsActive = true

	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			app.editConflictResponse(w, r)
		case errors.Is(err, models.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			return
		}
	}

	err = app.models.Tokens.DeleteAllTokenForUser(user.ID, models.ActivationScope)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"message": "user successfully updated"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) authenticateUserHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	user, err := app.models.Users.GetUserByEmail(input.Email)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.invalidUserCredintails(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	ok, err := user.Password.Matches(input.Password)
	if err != nil || !ok {
		app.serverErrorResponse(w, r, err)
		return
	}

	token, err := app.models.Tokens.New(24*time.Hour, user.ID, models.AuthenticationScope)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"token": token.Hash, "expiry": token.Expiry.Format(time.RFC822)}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
