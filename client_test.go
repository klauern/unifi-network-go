package unifi

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"testing"
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

	client, err := NewClient(baseURL, WithHTTPClient(httpClient))
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

func TestNewClient(t *testing.T) {
	tests := []struct {
		name    string
		baseURL string
		wantErr bool
	}{
		{
			name:    "valid URL",
			baseURL: "https://192.168.1.1",
			wantErr: false,
		},
		{
			name:    "invalid URL",
			baseURL: "://invalid",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewClient(tt.baseURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewClient() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && client == nil {
				t.Error("NewClient() returned nil client without error")
			}
		})
	}
}

func TestClient_do(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		expectedResp := PaginatedResponse{
			Offset:     0,
			Limit:      25,
			Count:      10,
			TotalCount: 1000,
			Data:       json.RawMessage(`[]`),
		}

		mock.response = mockResponse(http.StatusOK, expectedResp)

		var result PaginatedResponse
		err := client.do(ctx, http.MethodGet, "/v1/sites/test/hotspot/vouchers", nil, &result)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Limit != expectedResp.Limit {
			t.Errorf("expected limit %d, got %d", expectedResp.Limit, result.Limit)
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		apiError := Error{
			Status:      http.StatusNotFound,
			StatusName:  "Not Found",
			Message:     "Resource not found",
			RequestPath: "/v1/sites/test/hotspot/vouchers",
		}

		mock.response = mockResponse(http.StatusNotFound, apiError)

		err := client.do(ctx, http.MethodGet, "/v1/sites/test/hotspot/vouchers", nil, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		if apiErr, ok := err.(*Error); !ok {
			t.Errorf("expected *Error, got %T", err)
		} else if apiErr.Status != http.StatusNotFound {
			t.Errorf("expected status %d, got %d", http.StatusNotFound, apiErr.Status)
		}
	})

	t.Run("network error", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)
		mock.err = &url.Error{Op: "Get", URL: baseURL, Err: io.ErrUnexpectedEOF}

		err := client.do(ctx, http.MethodGet, "/v1/sites/test/hotspot/vouchers", nil, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestClient_ListHotspotVouchers(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()
	siteID := "default"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		expectedVouchers := []HotspotVoucher{
			{
				ID:                  "abc123",
				CreatedAt:           "2023-01-01T00:00:00Z",
				Name:                "Test Voucher",
				Code:                "WIFI123",
				TimeLimitMinutes:    60,
				DataUsageLimitMB:    1024,
				RxRateLimitKbps:     1024,
				TxRateLimitKbps:     512,
				AuthorizeGuestLimit: 2,
			},
		}

		mock.response = mockResponse(200, ListHotspotVouchersResponse{
			PaginatedResponse: PaginatedResponse{
				Count:      1,
				TotalCount: 1,
			},
			Data: expectedVouchers,
		})

		result, err := client.ListHotspotVouchers(ctx, siteID, &ListHotspotVouchersParams{
			Limit: 25,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Data) != 1 {
			t.Fatalf("expected 1 voucher, got %d", len(result.Data))
		}

		voucher := result.Data[0]
		if voucher.ID != expectedVouchers[0].ID {
			t.Errorf("expected voucher ID %s, got %s", expectedVouchers[0].ID, voucher.ID)
		}
		if voucher.Code != expectedVouchers[0].Code {
			t.Errorf("expected voucher code %s, got %s", expectedVouchers[0].Code, voucher.Code)
		}
	})

	t.Run("with pagination parameters", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(200, ListHotspotVouchersResponse{
			PaginatedResponse: PaginatedResponse{
				Offset:     50,
				Limit:      10,
				Count:      0,
				TotalCount: 100,
			},
			Data: []HotspotVoucher{},
		})

		params := &ListHotspotVouchersParams{
			Offset: 50,
			Limit:  10,
		}

		result, err := client.ListHotspotVouchers(ctx, siteID, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Offset != params.Offset {
			t.Errorf("expected offset %d, got %d", params.Offset, result.Offset)
		}
		if result.Limit != params.Limit {
			t.Errorf("expected limit %d, got %d", params.Limit, result.Limit)
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Site not found",
		})

		_, err := client.ListHotspotVouchers(ctx, "nonexistent", nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("expected *Error, got %T", err)
		}
	})
}

func TestClient_GetApplicationInfo(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		expectedInfo := ApplicationInfo{
			ApplicationVersion: "9.1.0",
		}

		mock.response = mockResponse(200, expectedInfo)

		result, err := client.GetApplicationInfo(ctx)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ApplicationVersion != expectedInfo.ApplicationVersion {
			t.Errorf("expected version %s, got %s", expectedInfo.ApplicationVersion, result.ApplicationVersion)
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(500, Error{
			Status:     500,
			StatusName: "Internal Server Error",
			Message:    "Server error",
		})

		_, err := client.GetApplicationInfo(ctx)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("expected *Error, got %T", err)
		}
	})
}
