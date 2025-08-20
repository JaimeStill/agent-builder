package main

import (
	"fmt"
	"net/http"
)

func (app *App) healthy(w http.ResponseWriter, r *http.Request) {
	health := "healthy"
	if !app.client.Healthy() {
		health = "sick"
	}

	fmt.Fprintf(w, "%s", health)
}

func (app *App) version(w http.ResponseWriter, r *http.Request) {
	response, err := app.client.Version()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	fmt.Fprintf(w, "%s", response.Version)
}
