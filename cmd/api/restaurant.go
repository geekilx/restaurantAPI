package main

import (
	"net/http"
	"strings"

	"github.com/geekilx/restaurantAPI/internal/models"
	"github.com/geekilx/restaurantAPI/internal/validator"
)

func (app *application) restaurantCreateHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string `json:"name"`
		Country     string `json:"country"`
		FullAddress string `json:"full_address"`
		Cuisine     string `json:"cuisine"`
		Status      string `json:"status"`
	}

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
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"restaurant": restaraunt}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
