package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) route() http.Handler {
	router := httprouter.New()

	// changing the default not found response in httprouter
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFoundResponse(w, r)
	})

	router.MethodNotAllowed = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		app.errorResponse(w, r, http.StatusMethodNotAllowed, "Method not allowed for this reousrce")
	})

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthCheck)
	router.HandlerFunc(http.MethodGet, "/v1/users/:id", app.requirePermissions("restaurant:read", app.userInformationHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users", app.createUserHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/users/:id", app.requirePermissions("restaurant:read", app.updateUserHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/resetpassword/:id", app.requirePermissions("restaurant:read", app.resetUserPasswordHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/users/:id", app.requirePermissions("restaurant:read", app.deleteUserHandler))
	router.HandlerFunc(http.MethodPost, "/v1/restaurants", app.requirePermissions("restaurant:write", app.restaurantCreateHandler))
	router.HandlerFunc(http.MethodGet, "/v1/restaurants", app.restaurantsListHandler)
	router.HandlerFunc(http.MethodGet, "/v1/restaurants/:id", app.showRestaurantHandler)
	router.HandlerFunc(http.MethodPatch, "/v1/restaurant/:id", app.requirePermissions("restaurant:write", app.restaurantUpdateHandler))
	router.HandlerFunc(http.MethodPost, "/v1/users/activate", app.userActivateHandler)
	router.HandlerFunc(http.MethodPost, "/v1/users/authenticate", app.authenticateUserHandler)
	router.HandlerFunc(http.MethodPost, "/v1/seller", app.createUserHandler)
	router.HandlerFunc(http.MethodGet, "/v1/category", app.allCategoryHandler)
	router.HandlerFunc(http.MethodPost, "/v1/category", app.requirePermissions("restaurant:write", app.createCategoryHandler))
	router.HandlerFunc(http.MethodGet, "/v1/restaurant/:id/categories", app.showSpecificRestaurantCategory)
	router.HandlerFunc(http.MethodPost, "/v1/category/:id/menu", app.requirePermissions("restaurant:write", app.createMenuHandler))
	router.HandlerFunc(http.MethodGet, "/v1/menus", app.menuListHandler)
	router.HandlerFunc(http.MethodGet, "/v1/category/:id", app.allMenuForCategoryHandler)

	return app.panicRecover(app.rateLimit(app.authenticate(router)))

}
