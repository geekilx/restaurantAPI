package main

import (
	"errors"
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
		return
	}

	user := app.getUserContext(r)

	if user.RestaurantID != nil {
		app.userHasRestaurant(w, r)
		return
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

	restaurantID, err := app.models.Restaurants.Insert(&restaraunt)
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

	user.RestaurantID = &restaurantID
	err = app.models.Users.Update(user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrRecordNotFound):
			app.noUserFound(w, r)
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

	var input struct {
		name string
		models.Filters
	}

	v := validator.New()

	qs := r.URL.Query()

	input.name = app.readString(qs, "name", "")
	input.Page = app.readInt(qs, "page", 1, v)
	input.PageSize = app.readInt(qs, "page_size", 20, v)
	input.Sort = app.readString(qs, "sort", "id")
	input.SortSafeList = []string{"id", "name", "country", "full_address", "cuisine", "status", "-id", "-name", "-country", "-full_address", "-cuisine", "-status"}

	restaraunts, metadata, err := app.models.Restaurants.GetAll(input.name, input.Filters)

	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	if restaraunts == nil {
		app.errorResponse(w, r, http.StatusOK, "there is no resaturant to be showen")
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"restaurants": restaraunts, "metadata": metadata}, nil)

}

func (app *application) restaurantUpdateHandler(w http.ResponseWriter, r *http.Request) {

	resID, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	var input struct {
		Name        string `json:"name"`
		Country     string `json:"country"`
		FullAddress string `json:"full_address"`
		Cuisine     string `json:"cuisine"`
		Status      string `json:"status"`
	}

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

// ShowRestaurant godoc
// @Summary      Get a single restaurant
// @Description  get string by ID
// @Tags         restaurants
// @Accept       json
// @Produce      json
// @Param        id   path      int  true  "Restaurant ID"
// @Success      200  {object}  jsFmt{menus=[]models.MenuWithCategoryName}
// @Failure      404  {object}  jsFmt
// @Failure      500  {object}  jsFmt
// @Router       /restaurants/{id} [get]
func (app *application) showRestaurantHandler(w http.ResponseWriter, r *http.Request) {

	restID, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	exists := app.models.Restaurants.CheckIfRestaurantExists(restID)
	if !exists {
		app.noRestaurantFound(w, r)
		return
	}

	menus, err := app.models.Menu.GetRestaurantMenus(restID)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if menus == nil {
		app.noMenuAvailable(w, r)
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"menus": menus}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
