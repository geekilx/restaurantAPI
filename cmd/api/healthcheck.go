package main

import (
	"net/http"
)

func (app *application) healthCheck(w http.ResponseWriter, r *http.Request) {
	msg := map[string]any{
		"status":  "available",
		"version": Version,
	}

	app.writeJSON(w, r, http.StatusOK, jsFmt{"message": msg}, nil)

}
