package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/JaimeStill/agent-builder/pkg/ollama"
)

func main() {
	addr := flag.String("addr", ":5000", "HTTP network address")
	ollamaEndpoint := flag.String("ollama", "http://localhost:11434", "Ollama endpoint")
	flag.Parse()

	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))
	options := ollama.DefaultOptions(*ollamaEndpoint)

	app := &App{
		logger: logger,
		client: ollama.NewClient(&options),
	}

	srv := &http.Server{
		Addr:         *addr,
		Handler:      app.routes(),
		ErrorLog:     slog.NewLogLogger(logger.Handler(), slog.LevelError),
		ReadTimeout:  1 * time.Minute,
		WriteTimeout: 60 * time.Minute,
		IdleTimeout:  2 * time.Minute,
	}

	logger.Info("starting server", "addr", srv.Addr)
	err := srv.ListenAndServe()
	logger.Error(err.Error())
	os.Exit(1)
}
