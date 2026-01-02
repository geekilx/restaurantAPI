package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/geekilx/restaurantAPI/internal/models"
	"github.com/geekilx/restaurantAPI/internal/validator"
)

var input struct {
	Name        string `json:"name"`
	Country     string `json:"country"`
	FullAddress string `json:"full_address"`
	Cuisine     string `json:"cuisine"`
	Status      string `json:"status"`
}

func (app *application) restaurantCreateHandler(w http.ResponseWriter, r *http.Request) {

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)

	}

	restaraunt := models.Restaurant{
		Name:        input.Name,
		Country:     input.Country,
		FullAddress: input.FullAddress,
		Cuisine:     input.Cuisine,
		Status:      strings.ToLower(input.Status),
	}

	v := validator.New()

	if models.ValidateRestaurant(v, restaraunt); !v.Valid() {
		app.failedValidationResponse(w, r, v)
		return
	}

	err = app.models.Restaurants.Insert(&restaraunt)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateRestaurantName):
			v.AddError("restaurant name", "the restaurant name was already taken")
			app.failedValidationResponse(w, r, v)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, r, http.StatusCreated, jsFmt{"restaurant": restaraunt}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) restaurantsListHandler(w http.ResponseWriter, r *http.Request) {
	restaraunts, err := app.models.Restaurants.GetAll()

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if restaraunts == nil {
		app.errorResponse(w, r, http.StatusOK, "there is no resaturant to be showen")
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"restaurants": restaraunts}, nil)

}

func (app *application) restaurantUpdateHandler(w http.ResponseWriter, r *http.Request) {

	resID, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	// country      | character varying(50)    |           | not null |
	// full_address | text                     |           | not null |
	// cuisine      | character varying(50)    |           | not null |
	// status       | character varying(50)    |           | not null |

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	restaurant, err := app.models.Restaurants.Get(resID)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRestaurantNotFound):
			app.noRestaurantFound(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	if input.Name != "" {
		restaurant.Name = input.Name
	}
	if input.Country != "" {
		restaurant.Country = input.Country
	}
	if input.FullAddress != "" {
		restaurant.FullAddress = input.FullAddress
	}
	if input.Cuisine != "" {
		restaurant.Cuisine = input.Cuisine
	}
	if input.Status != "" {
		restaurant.Status = input.Status
	}

	v := validator.New()

	if models.ValidateRestaurant(v, *restaurant); !v.Valid() {
		app.failedValidationResponse(w, r, v)
		return
	}

	err = app.models.Restaurants.Update(resID, *restaurant)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateRestaurantName):
			v.AddError("restaurant name", "the resaturant name already exists")
			app.failedValidationResponse(w, r, v)
		case errors.Is(err, models.ErrRestaurantNotFound):
			app.noRestaurantFound(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"message": "The resaturant was updated successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *application) restaurantDeleteHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.models.Restaurants.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRestaurantNotFound):
			app.noRestaurantFound(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"message": "the resaturant was deleted successfully"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
