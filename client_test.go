package unifi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"testing"
)

func TestNewClient(t *testing.T) {
	t.Run("valid URL", func(t *testing.T) {
		_, err := NewClient(
			"https://192.168.1.1",
			WithAPIKey("test-api-key"),
		)
		if err != nil {
			t.Errorf("NewClient() error = %v, wantErr false", err)
		}
	})

	t.Run("invalid URL", func(t *testing.T) {
		_, err := NewClient(
			"://invalid",
			WithAPIKey("test-api-key"),
		)
		if err == nil {
			t.Error("NewClient() error = nil, wantErr true")
		}
	})
}

func TestClient_do(t *testing.T) {
	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		expectedResponse := struct {
			Message string `json:"message"`
		}{
			Message: "success",
		}

		mock.response = mockResponse(200, expectedResponse)

		var result struct {
			Message string `json:"message"`
		}

		err := client.do(context.Background(), http.MethodGet, "/test", nil, &result)
		if err != nil {
			t.Errorf("do() error = %v", err)
		}

		if result.Message != expectedResponse.Message {
			t.Errorf("do() got = %v, want %v", result.Message, expectedResponse.Message)
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(400, Error{
			Status:     400,
			StatusName: "Bad Request",
			Message:    "Invalid parameters",
		})

		err := client.do(context.Background(), http.MethodGet, "/test", nil, nil)
		if err == nil {
			t.Error("do() error = nil, wantErr true")
		}

		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("do() error type = %T, want *Error", err)
		}
	})

	t.Run("network error", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		expectedErr := fmt.Errorf("network error")
		mock.err = expectedErr

		err := client.do(context.Background(), http.MethodGet, "/test", nil, nil)
		if err == nil {
			t.Error("do() error = nil, wantErr true")
		}

		if !strings.Contains(err.Error(), expectedErr.Error()) {
			t.Errorf("do() error = %v, want %v", err, expectedErr)
		}
	})
}

func TestClient_ListHotspotVouchers(t *testing.T) {
	ctx := context.Background()

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

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
				Offset:     0,
				Limit:      25,
			},
			Data: expectedVouchers,
		})

		result, err := client.ListHotspotVouchers(ctx, testSiteID, &ListHotspotVouchersParams{
			Limit: 25,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertPaginatedResponse(t, result.PaginatedResponse, PaginatedResponse{
			Count:      1,
			TotalCount: 1,
			Offset:     0,
			Limit:      25,
		})

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
		client, mock := newTestClient(t, testBaseURL)

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

		result, err := client.ListHotspotVouchers(ctx, testSiteID, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertPaginatedResponse(t, result.PaginatedResponse, PaginatedResponse{
			Offset:     50,
			Limit:      10,
			Count:      0,
			TotalCount: 100,
		})
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Site not found",
		})

		_, err := client.ListHotspotVouchers(ctx, "nonexistent", nil)
		assertErrorResponse(t, err, 404, "Site not found")
	})
}

func TestClient_GetApplicationInfo(t *testing.T) {
	ctx := context.Background()

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

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
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(500, Error{
			Status:     500,
			StatusName: "Internal Server Error",
			Message:    "Server error",
		})

		_, err := client.GetApplicationInfo(ctx)
		assertErrorResponse(t, err, 500, "Server error")
	})
}
