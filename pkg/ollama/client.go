package ollama

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
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

func (client *Client) Pull(ctx context.Context, req *PullRequest, handler func(*PullResponse) error) error {
	body, err := json.Marshal(req)
	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/api/pull", client.Endpoint())
	request, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json")

	pullClient := &http.Client{
		Transport: client.http.Transport,
	}

	response, err := pullClient.Do(request)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		var errBody bytes.Buffer
		io.Copy(&errBody, response.Body)
		return fmt.Errorf("pull failed: %s", errBody.String())
	}

	if req.GetStream() {
		return client.handleStreamingResponse(response.Body, handler)
	} else {
		return client.handleSingleResponse(response.Body, handler)
	}
}

func (client *Client) PS() (*ModelResponse, error) {
	url := fmt.Sprintf("%s/api/ps", client.Endpoint())
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.http.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var model ModelResponse
	if err = json.NewDecoder(response.Body).Decode(&model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (client *Client) List() (*ModelResponse, error) {
	url := fmt.Sprintf("%s/api/tags", client.Endpoint())
	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	response, err := client.http.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var model ModelResponse
	if err = json.NewDecoder(response.Body).Decode(&model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (client *Client) Show(req *ShowRequest) (*ShowResponse, error) {
	body, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	url := fmt.Sprintf("%s/api/show", client.Endpoint())
	request, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	request.Header.Set("Content-Type", "application/json")

	response, err := client.http.Do(request)
	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	var model ShowResponse
	if err = json.NewDecoder(response.Body).Decode(&model); err != nil {
		return nil, err
	}

	return &model, nil
}

func (client *Client) handleStreamingResponse(body io.Reader, handler func(*PullResponse) error) error {
	decoder := json.NewDecoder(body)

	for {
		var response PullResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				return nil
			}
			return fmt.Errorf("failed to decode response: %w", err)
		}

		if err := handler(&response); err != nil {
			return err
		}

		if response.Status == "success" {
			break
		}
	}

	return nil
}

func (client *Client) handleSingleResponse(body io.Reader, handler func(*PullResponse) error) error {
	var response PullResponse
	if err := json.NewDecoder(body).Decode(&response); err != nil {
		return err
	}

	return handler(&response)
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
