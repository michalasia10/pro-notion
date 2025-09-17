package notion

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	authTypeBearer = "Bearer"
	authTypeBasic  = "Basic"
)

// Client represents a Notion API client
type Client struct {
	baseURL    string
	httpClient *http.Client
	apiVersion string
}

// ClientOption represents a configuration option for the Client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(client *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = client
	}
}

// WithAPIVersion sets a custom API version
func WithAPIVersion(version string) ClientOption {
	return func(c *Client) {
		c.apiVersion = version
	}
}

// NewClient creates a new Notion API client
func NewClient(opts ...ClientOption) *Client {
	client := &Client{
		baseURL:    "https://api.notion.com/v1",
		apiVersion: "2022-06-28",
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Jar:     nil,
		},
	}

	for _, opt := range opts {
		opt(client)
	}

	return client
}

// makeRequest performs an HTTP request to the Notion API
func (c *Client) makeRequest(method, path string, body interface{}, token string, authType string) (*http.Response, error) {
	var reqBody io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonBody)
	}

	url := fmt.Sprintf("%s%s", c.baseURL, path)
	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Notion-Version", c.apiVersion)

	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("%s %s", authType, token))
	}
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return resp, fmt.Errorf("failed to execute request: %w", err)
	}

	return resp, nil
}

// handleResponse handles the API response
func (c *Client) handleResponse(resp *http.Response, result interface{}) error {
	if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusBadRequest {
		var apiErr APIError
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response body: %w", err)
		}

		defer resp.Body.Close()

		if err := json.Unmarshal(bodyBytes, &apiErr); err != nil {
			return fmt.Errorf("failed to unmarshal error response (status: %s): %s", resp.Status, string(bodyBytes))
		}
		return &apiErr
	}

	defer resp.Body.Close()
	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}
