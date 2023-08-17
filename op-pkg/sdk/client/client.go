package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/grafana/grafana/op-pkg/sdk/client/interceptor"
)

var (
	// ErrNotFound is returned when there is no such entity by given ID
	ErrNotFound = errors.New("not found")
)

// Client is a customizable http client
type Client struct {
	baseURL    string
	headers    map[string]string
	httpClient *http.Client
}

// New creates new Client
func New(baseURL string, options ...Option) *Client {
	client := &Client{
		baseURL:    baseURL,
		headers:    make(map[string]string),
		httpClient: http.DefaultClient,
	}
	for _, opt := range options {
		opt(client)
	}
	return client
}

// Get performs GET request
func (c *Client) Get(ctx context.Context, endpoint string, interceptors ...interceptor.Interceptor) ([]byte, error) {
	return c.processRequest(ctx, http.MethodGet, endpoint, http.NoBody, interceptors...)
}

// Delete performs DELETE request
func (c *Client) Delete(ctx context.Context, endpoint string, interceptors ...interceptor.Interceptor) ([]byte, error) {
	return c.processRequest(ctx, http.MethodDelete, endpoint, http.NoBody, interceptors...)
}

// Post performs POST request
func (c *Client) Post(ctx context.Context, endpoint string, r io.Reader, interceptors ...interceptor.Interceptor) ([]byte, error) {
	return c.processRequest(ctx, http.MethodPost, endpoint, r, interceptors...)
}

func (c *Client) processRequest(ctx context.Context, method, endpoint string, r io.Reader, ics ...interceptor.Interceptor) ([]byte, error) {
	u, err := url.JoinPath(c.baseURL, endpoint)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequestWithContext(ctx, method, u, r)
	if err != nil {
		return nil, fmt.Errorf("invalid request parameters: %w", err)
	}
	for k, v := range c.headers {
		request.Header.Set(k, v)
	}
	var do = interceptor.Doer(c.httpClient.Do)
	for _, ic := range ics {
		do = ic(do)
	}
	response, doErr := do(request)
	if doErr != nil {
		return nil, doErr
	}
	defer response.Body.Close()
	data, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}
	if isValidStatusCode(response.StatusCode) {
		return data, nil
	}
	if message := extractErrorMessage(data); message != "" {
		return nil, fmt.Errorf("invalid status: %s, message: %s", response.Status, message)
	}
	return nil, fmt.Errorf("invalid status: %s", response.Status)
}

func isValidStatusCode(statusCode int) bool {
	return statusCode == http.StatusOK || statusCode == http.StatusCreated
}

func extractErrorMessage(data []byte) (msg string) {
	return string(data)
}
