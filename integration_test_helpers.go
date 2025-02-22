package unifi

import (
	"crypto/tls"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

// integrationTestConfig holds configuration for integration tests
type integrationTestConfig struct {
	BaseURL  string
	APIKey   string
	Insecure bool
}

// skipIfNotIntegration skips the test if integration tests are not enabled
func skipIfNotIntegration(t *testing.T) {
	t.Helper()
	if os.Getenv("UNIFI_INTEGRATION_TEST") != "1" {
		t.Skip("Skipping integration test. Set UNIFI_INTEGRATION_TEST=1 to run")
	}
}

// loadIntegrationConfig loads configuration for integration tests
func loadIntegrationConfig(t *testing.T) *integrationTestConfig {
	t.Helper()
	skipIfNotIntegration(t)

	baseURL := os.Getenv("UNIFI_BASE_URL")
	if baseURL == "" {
		t.Fatal("UNIFI_BASE_URL environment variable is required for integration tests")
	}

	// Parse the URL to check if it has a port
	parsedURL, err := url.Parse(baseURL)
	if err != nil {
		t.Fatalf("failed to parse base URL: %v", err)
	}

	// Add port 8443 if no port is specified
	if parsedURL.Port() == "" {
		host := parsedURL.Host
		if strings.Contains(host, ":") {
			host = host[:strings.Index(host, ":")]
		}
		parsedURL.Host = fmt.Sprintf("%s:8443", host)
		baseURL = parsedURL.String()
	}

	// Ensure the base URL has the correct structure for UniFi Network API
	if !strings.Contains(baseURL, "/proxy/network") {
		baseURL = strings.TrimSuffix(baseURL, "/") + "/proxy/network"
	}

	apiKey := os.Getenv("UNIFI_API_KEY")
	if apiKey == "" {
		t.Fatal("UNIFI_API_KEY environment variable is required for integration tests")
	}

	// Default to insecure for integration tests since many UniFi controllers use self-signed certs
	insecure := true
	if os.Getenv("UNIFI_INTEGRATION_SECURE") == "1" {
		insecure = false
	}

	return &integrationTestConfig{
		BaseURL:  baseURL,
		APIKey:   apiKey,
		Insecure: insecure,
	}
}

// newIntegrationTestClient creates a new client for integration testing
func newIntegrationTestClient(t *testing.T) *Client {
	t.Helper()
	config := loadIntegrationConfig(t)

	// Create a custom HTTP client that allows insecure TLS if needed
	httpClient := &http.Client{}
	if config.Insecure {
		httpClient.Transport = &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}
	}

	client, err := NewClient(
		config.BaseURL,
		WithHTTPClient(httpClient),
		WithAPIKey(config.APIKey),
	)
	if err != nil {
		t.Fatalf("failed to create integration test client: %v", err)
	}

	return client
}
