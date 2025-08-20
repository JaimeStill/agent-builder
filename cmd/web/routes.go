package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *App) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /ollama/health", app.health)

	standard := alice.New(
		app.recoverPanic,
		app.logRequest,
		commonHeaders,
	)

	return standard.Then(mux)
}
