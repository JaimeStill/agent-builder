package main

import (
	"encoding/json"
	"fmt"
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

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *App) version(w http.ResponseWriter, r *http.Request) {
	response, err := app.client.Version()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *App) pull(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		app.clientError(w, http.StatusMethodNotAllowed)
		return
	}

	var pullReq ollama.PullRequest
	if err := json.NewDecoder(r.Body).Decode(&pullReq); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	if pullReq.GetStream() {
		w.Header().Set("Content-Type", "application/x-ndjson")
	} else {
		w.Header().Set("Content-Type", "application.json")
	}

	flusher, canFlush := w.(http.Flusher)

	err := app.client.Pull(r.Context(), &pullReq, func(resp *ollama.PullResponse) error {
		data, err := json.Marshal(resp)
		if err != nil {
			return err
		}

		if pullReq.GetStream() {
			fmt.Fprintf(w, "%s\n", data)
			if canFlush {
				flusher.Flush()
			}
		} else {
			w.Write(data)
		}

		return nil
	})

	if err != nil {
		if pullReq.GetStream() {
			errResp := map[string]string{"error": err.Error()}
			data, _ := json.Marshal(errResp)
			fmt.Fprintf(w, "%s\n", data)
			if canFlush {
				flusher.Flush()
			}
		} else {
			app.serverError(w, r, err)
		}
	}
}

func (app *App) ps(w http.ResponseWriter, r *http.Request) {
	response, err := app.client.PS()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.serverError(w, r, err)
	}
}

func (app *App) list(w http.ResponseWriter, r *http.Request) {
	response, err := app.client.List()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
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

	response, err := app.client.Show(&showReq)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		app.serverError(w, r, err)
	}
}
