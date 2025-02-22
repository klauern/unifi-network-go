package unifi

import (
	"crypto/tls"
	"net/http"
	"os"
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
