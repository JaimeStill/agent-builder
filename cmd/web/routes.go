package main

import (
	"net/http"

	"github.com/justinas/alice"
)

func (app *App) routes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", app.healthy)
	mux.HandleFunc("GET /ollama/version", app.version)
	mux.HandleFunc("POST /ollama/pull", app.pull)
	mux.HandleFunc("GET /ollama/ps", app.ps)
	mux.HandleFunc("GET /ollama/list", app.list)
	mux.HandleFunc("POST /ollama/show", app.show)

	standard := alice.New(
		app.recoverPanic,
		app.logRequest,
		commonHeaders,
	)

	return standard.Then(mux)
}
