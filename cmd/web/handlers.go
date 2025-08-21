package main

import (
	"encoding/json"
	"net/http"

	"github.com/JaimeStill/agent-builder/pkg/ollama"
)

func (app *App) healthy(w http.ResponseWriter, r *http.Request) {
	health := "healthy"
	if !app.client.Healthy() {
		health = "sick"
	}

	response := struct {
		Status string `json:"status"`
	}{
		Status: health,
	}

	if err := app.writeJSON(w, http.StatusOK, response, nil); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *App) version(w http.ResponseWriter, r *http.Request) {
	response, err := app.client.Version()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, response, nil); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *App) pull(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var pullReq ollama.PullRequest
	defer r.Body.Close()
	if err := json.NewDecoder(r.Body).Decode(&pullReq); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if !pullReq.GetStream() {
		w.Header().Set("Content-Type", "application/json")
		var responses []*ollama.PullResponse

		err := app.client.Pull(r.Context(), &pullReq, func(resp *ollama.PullResponse) error {
			responses = append(responses, resp)
			return nil
		})

		if err != nil {
			app.serverError(w, r, err)
			return
		}

		app.writeJSON(w, http.StatusOK, responses, nil)
		return
	}

	w.Header().Set("Content-Type", "application/x-ndjson")
	w.WriteHeader(http.StatusOK)

	flusher, canFlush := w.(http.Flusher)
	encoder := json.NewEncoder(w)

	err := app.client.Pull(r.Context(), &pullReq, func(resp *ollama.PullResponse) error {
		if err := encoder.Encode(resp); err != nil {
			return err
		}
		if canFlush {
			flusher.Flush()
		}
		return nil
	})

	if err != nil {
		encoder.Encode(map[string]string{"error": err.Error()})
		if canFlush {
			flusher.Flush()
		}
	}
}

func (app *App) ps(w http.ResponseWriter, r *http.Request) {
	response, err := app.client.PS()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, response, nil); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *App) list(w http.ResponseWriter, r *http.Request) {
	response, err := app.client.List()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, response, nil); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *App) show(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var showReq ollama.ShowRequest
	if err := json.NewDecoder(r.Body).Decode(&showReq); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	response, err := app.client.Show(&showReq)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	if err := app.writeJSON(w, http.StatusOK, response, nil); err != nil {
		app.serverError(w, r, err)
	}
}
