package ollama

import "time"

type Options struct {
	Endpoint         string        `json:"endpoint"`
	Timeout          time.Duration `json:"timeout"`
	MaxRetries       int           `json:"max_retries"`
	RetryBackoffBase time.Duration `json:"retry_backoff_base"`
	MaxConnections   int           `json:"max_connections"`
	IdleTimeout      time.Duration `json:"idle_timeout"`
}

func DefaultOptions(endpoint string) Options {
	return Options{
		Endpoint:         endpoint,
		Timeout:          60 * time.Second,
		MaxRetries:       3,
		RetryBackoffBase: time.Second,
		MaxConnections:   10,
		IdleTimeout:      90 * time.Second,
	}
}
