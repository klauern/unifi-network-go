package unifi

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
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
	apiKey     string
	insecure   bool
	logger     *slog.Logger
}

// ClientOption allows for customizing the client
type ClientOption func(*Client)

// WithHTTPClient sets a custom HTTP client
func WithHTTPClient(httpClient *http.Client) ClientOption {
	return func(c *Client) {
		c.httpClient = httpClient
	}
}

// WithAPIKey sets the API key for authentication
func WithAPIKey(apiKey string) ClientOption {
	return func(c *Client) {
		c.apiKey = apiKey
	}
}

// WithInsecure sets whether to skip TLS certificate verification
func WithInsecure(insecure bool) ClientOption {
	return func(c *Client) {
		c.insecure = insecure
	}
}

// WithLogger sets a custom logger for the client
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// NewClient creates a new UniFi Network API client
func NewClient(baseURL string, options ...ClientOption) (*Client, error) {
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		return nil, fmt.Errorf("invalid base URL: %w", err)
	}

	// Ensure the base path includes the API prefix
	// First, trim any existing proxy/network/integration prefix to avoid doubles
	trimmedPath := strings.TrimPrefix(parsedURL.Path, "/proxy/network/integration")
	trimmedPath = strings.TrimPrefix(trimmedPath, "proxy/network/integration")
	parsedURL.Path = path.Join("/proxy/network/integration", trimmedPath)

	// Create default logger
	logLevel := new(slog.LevelVar)
	if os.Getenv("DEBUG") != "" {
		logLevel.Set(slog.LevelDebug)
	} else {
		logLevel.Set(slog.LevelInfo)
	}
	defaultLogger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: logLevel,
	}))

	client := &Client{
		baseURL:    parsedURL,
		httpClient: http.DefaultClient,
		logger:     defaultLogger,
	}

	for _, opt := range options {
		opt(client)
	}

	if client.apiKey == "" {
		return nil, fmt.Errorf("API key is required")
	}

	// Configure TLS if insecure is set
	if client.insecure {
		transport := http.DefaultTransport.(*http.Transport).Clone()
		transport.TLSClientConfig = &tls.Config{
			InsecureSkipVerify: true,
		}
		client.httpClient = &http.Client{
			Transport: transport,
		}
	}

	client.logger.Debug("Created UniFi Network client",
		"base_url", client.baseURL.String(),
		"insecure", client.insecure)

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

	// Split the path and query if present
	pathParts := strings.Split(urlPath, "?")
	u.Path = path.Join(u.Path, pathParts[0])

	// Add query parameters if they exist
	if len(pathParts) > 1 {
		u.RawQuery = pathParts[1]
	}

	c.logger.Debug("Constructing request",
		"base_path", c.baseURL.Path,
		"url_path", urlPath,
		"final_path", u.Path,
		"query_params", u.RawQuery,
		"final_url", u.String())

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
		bodyReader = bytes.NewReader(jsonBody)
		c.logger.Debug("Request body", "body", string(jsonBody))
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("X-API-KEY", c.apiKey)

	c.logger.Debug("Making request",
		"method", method,
		"url", u.String(),
		"headers", req.Header)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read the entire response body for debugging
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	c.logger.Debug("Received response",
		"status", resp.Status,
		"body_length", len(respBody))

	if resp.StatusCode >= 400 {
		var apiErr Error
		if err := json.Unmarshal(respBody, &apiErr); err != nil {
			// If we can't decode the error response, return the raw response
			return fmt.Errorf("API error (status %d): %s", resp.StatusCode, string(respBody))
		}
		return &apiErr
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("failed to decode response: %w\nResponse body: %s", err, string(respBody))
		}
	}

	return nil
}

// Error implements the error interface for UniFi API errors
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s (status: %d, request: %s, id: %s)",
		e.StatusName, e.Message, e.Status, e.RequestPath, e.RequestID)
}
