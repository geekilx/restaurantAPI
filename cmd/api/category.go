package main

import (
	"fmt"
	"net/http"

	"github.com/geekilx/restaurantAPI/internal/models"
	"github.com/geekilx/restaurantAPI/internal/validator"
)

func (app *application) createCategoryHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name string `json:"name"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	v := validator.New()

	user := app.getUserContext(r)
	category := models.Category{
		RestaurantID: *user.RestaurantID,
		Name:         input.Name,
	}

	ok := app.models.Categories.CategoryExists(category.Name, *user.RestaurantID)
	if ok {
		v.AddError("name", "this restaurant has already created this category")
		app.failedValidationResponse(w, r, v)
		return
	}

	err = app.models.Categories.Insert(&category)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"category": category}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) allCategoryHandler(w http.ResponseWriter, r *http.Request) {

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
	input.SortSafeList = []string{"id", "name", "restaurant_id", "-id", "-name", "-restaurant_id"}

	categories, metadata, err := app.models.Categories.GetAll(input.name, input.Filters)
	if err != nil {
		app.noCategoryIsAvailable(w, r)
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"Categories": categories, "metadata": metadata}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}

func (app *application) showSpecificRestaurantCategory(w http.ResponseWriter, r *http.Request) {
	restID, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	fmt.Println(restID)

	catgeories, err := app.models.Categories.GetAllForRestaurant(restID)
	if err != nil || catgeories == nil {
		app.noCategoryIsAvailable(w, r)
		return
	}

	err = app.writeJSON(w, r, http.StatusOK, jsFmt{"categories": catgeories}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) allMenuForCategoryHandler(w http.ResponseWriter, r *http.Request) {

	catID, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	menus, err := app.models.Menu.GetAllMenuForCategory(catID)
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
