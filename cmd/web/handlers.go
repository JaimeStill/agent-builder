package main

import (
	"fmt"
	"net/http"
)

func (app *App) health(w http.ResponseWriter, r *http.Request) {
	health, err := app.client.Health()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	fmt.Fprintf(w, "%v", health)
}
