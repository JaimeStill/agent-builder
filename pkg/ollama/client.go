package ollama

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Client struct {
	options    *Options
	http       *http.Client
	mu         sync.RWMutex
	healthy    bool
	lastHealth time.Time
}

func NewClient(options *Options) *Client {
	client := &Client{
		options: options,
		healthy: false,
		http: &http.Client{
			Timeout: options.Timeout,
			Transport: &http.Transport{
				MaxIdleConns:        options.MaxConnections,
				MaxIdleConnsPerHost: options.MaxConnections,
				IdleConnTimeout:     options.IdleTimeout,
			},
		},
	}

	client.performHealthCheck()

	return client
}

func (client *Client) Endpoint() string {
	return client.options.Endpoint
}

func (client *Client) Close() {
	if transport, ok := client.http.Transport.(*http.Transport); ok {
		transport.CloseIdleConnections()
	}
}

func (client *Client) Health() (*HealthResponse, error) {
	url := fmt.Sprintf("%s/api/version", client.Endpoint())
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.http.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var healthResponse HealthResponse
	if err = json.NewDecoder(response.Body).Decode(&healthResponse); err != nil {
		return nil, err
	}

	return &healthResponse, nil
}

func (client *Client) performHealthCheck() {
	_, err := client.Health()
	client.setHealthy(err == nil)
}

func (client *Client) setHealthy(healthy bool) {
	client.mu.Lock()
	defer client.mu.Unlock()

	client.healthy = healthy
	client.lastHealth = time.Now()
}
