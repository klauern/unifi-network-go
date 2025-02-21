package unifi

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path"
)

type Error struct {
	Status      int    `json:"statusCode"`
	StatusName  string `json:"statusName"`
	Message     string `json:"message"`
	Timestamp   string `json:"timestamp"`
	RequestPath string `json:"requestPath"`
	RequestID   string `json:"requestId"`
}

// Client represents a UniFi Network API client
type Client struct {
	baseURL    *url.URL
	httpClient *http.Client
}

// ClientOption allows for customizing the client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// NewClient creates a new UniFi Network API client
func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	client := &Client{
		baseURL:    parsedURL,
		httpClient: http.DefaultClient,
	}

	for _, opt := range options {
		opt(client)
	}

	return client, nil
}

// PaginatedResponse represents a paginated API response
type PaginatedResponse struct {
	Offset     int             `json:"offset"`
	Limit      int             `json:"limit"`
	Count      int             `json:"count"`
	TotalCount int             `json:"totalCount"`
	Data       json.RawMessage `json:"data"`
}

// ApplicationInfo represents the UniFi Network application information
type ApplicationInfo struct {
	ApplicationVersion string `json:"applicationVersion"` // Version of the UniFi Network application
}

// GetApplicationInfo retrieves generic information about the Network application
func (c *Client) GetApplicationInfo(ctx context.Context) (*ApplicationInfo, error) {
	var response ApplicationInfo
	err := c.do(ctx, http.MethodGet, "/v1/info", nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get application info: %w", err)
	}

	return &response, nil
}

func (c *Client) do(ctx context.Context, method, urlPath string, body interface{}, result interface{}) error {
	u := *c.baseURL
	u.Path = path.Join(u.Path, urlPath)

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		var apiErr Error
		if err := json.NewDecoder(resp.Body).Decode(&apiErr); err != nil {
			return fmt.Errorf("failed to decode error response: %w", err)
		}
		return &apiErr
	}

	if result != nil {
		if err := json.NewDecoder(resp.Body).Decode(result); err != nil {
			return fmt.Errorf("failed to decode response: %w", err)
		}
	}

	return nil
}

// Error implements the error interface for UniFi API errors
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s (status: %d, request: %s, id: %s)",
		e.StatusName, e.Message, e.Status, e.RequestPath, e.RequestID)
}
