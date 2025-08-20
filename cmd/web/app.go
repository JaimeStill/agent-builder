package main

import (
	"log/slog"
	"net/http"

	"github.com/JaimeStill/agent-builder/pkg/ollama"
)

type App struct {
	logger *slog.Logger
	client *ollama.Client
}

func (app *App) serverError(w http.ResponseWriter, r *http.Request, err error) {
	var (
		method = r.Method
		uri    = r.URL.RequestURI()
	)

	app.logger.Error(err.Error(), "method", method, "uri", uri)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

func (app *App) clientError(w http.ResponseWriter, status int) {
	http.Error(w, http.StatusText(status), status)
}
