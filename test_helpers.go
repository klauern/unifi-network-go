package unifi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"testing"
)

const (
	testBaseURL = "https://192.168.1.1"
	testSiteID  = "default"
)

// mockTransport implements http.RoundTripper for testing
type mockTransport struct {
	response *http.Response
	err      error
}

func (t *mockTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return t.response, t.err
}

// newTestClient creates a client with a mock transport for testing
func newTestClient(t *testing.T, baseURL string) (*Client, *mockTransport) {
	t.Helper()
	mock := &mockTransport{}
	httpClient := &http.Client{Transport: mock}

	client, err := NewClient(
		baseURL,
		WithHTTPClient(httpClient),
		WithAPIKey("test-api-key"),
	)
	if err != nil {
		t.Fatalf("failed to create test client: %v", err)
	}

	return client, mock
}

// mockResponse creates a mock HTTP response with the given status code and body
func mockResponse(statusCode int, body interface{}) *http.Response {
	var bodyReader io.ReadCloser
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			panic(err)
		}
		bodyReader = io.NopCloser(bytes.NewReader(b))
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       bodyReader,
	}
}

// assertPaginatedResponse validates common pagination fields
func assertPaginatedResponse(t *testing.T, got, want PaginatedResponse) {
	t.Helper()
	if got.Count != want.Count {
		t.Errorf("expected count %d, got %d", want.Count, got.Count)
	}
	if got.TotalCount != want.TotalCount {
		t.Errorf("expected total count %d, got %d", want.TotalCount, got.TotalCount)
	}
	if got.Offset != want.Offset {
		t.Errorf("expected offset %d, got %d", want.Offset, got.Offset)
	}
	if got.Limit != want.Limit {
		t.Errorf("expected limit %d, got %d", want.Limit, got.Limit)
	}
}

// assertErrorResponse validates error responses
func assertErrorResponse(t *testing.T, err error, wantStatus int, wantMessage string) {
	t.Helper()
	var apiErr *Error
	if !errors.As(err, &apiErr) {
		t.Errorf("expected *Error, got %T", err)
		return
	}
	if apiErr.Status != wantStatus {
		t.Errorf("expected status %d, got %d", wantStatus, apiErr.Status)
	}
	if apiErr.Message != wantMessage {
		t.Errorf("expected message %q, got %q", wantMessage, apiErr.Message)
	}
}
