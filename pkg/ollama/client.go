package ollama

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type Client struct {
	options *Options
	http    *http.Client
	mu      sync.RWMutex
	healthy bool
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

func (client *Client) Healthy() bool {
	return client.healthy
}

func (client *Client) Version() (*VersionResponse, error) {
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

	var version VersionResponse
	if err = json.NewDecoder(response.Body).Decode(&version); err != nil {
		return nil, err
	}

	return &version, nil
}

func (client *Client) performHealthCheck() {
	_, err := client.Version()
	client.setHealthy(err == nil)
}

func (client *Client) setHealthy(healthy bool) {
	client.mu.Lock()
	defer client.mu.Unlock()

	client.healthy = healthy
}
