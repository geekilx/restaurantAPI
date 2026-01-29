package main

import (
	"net/http"

	"github.com/geekilx/restaurantAPI/internal/models"
	"github.com/geekilx/restaurantAPI/internal/validator"
)

func (app *application) createMenuHandler(w http.ResponseWriter, r *http.Request) {

	categoryID, err := app.readIDParam(r)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	exists := app.models.Categories.CheckIfExists(categoryID)
	if !exists {
		app.noCategoryIsAvailable(w, r)
		return
	}

	var input struct {
		Name        string  `json:"name"`
		Description string  `json:"description"`
		PriceCent   float32 `json:"price_cent"`
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	menu := models.Menu{
		CategoryID:  categoryID,
		Name:        input.Name,
		Description: input.Description,
		PriceCent:   input.PriceCent,
	}

	err = app.models.Menu.Insert(&menu)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, r, http.StatusCreated, jsFmt{"menu": menu}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
func (app *application) menuListHandler(w http.ResponseWriter, r *http.Request) {

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
	input.SortSafeList = []string{"id", "name", "description", "price_cent", "is_available", "-id", "-name", "-description", "-price_cent", "-is_available"}

	menus, err := app.models.Menu.GetAll(input.name, input.Filters)
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
