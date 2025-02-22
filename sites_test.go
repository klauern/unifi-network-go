package unifi

import (
	"context"
	"testing"
)

func TestClient_ListSites(t *testing.T) {
	ctx := context.Background()

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		expectedSites := []Site{
			{
				ID:   "default",
				Name: "Default",
			},
		}

		mock.response = mockResponse(200, ListSitesResponse{
			Offset:     0,
			Limit:      25,
			Count:      1,
			TotalCount: 1,
			Data:       expectedSites,
		})

		result, err := client.ListSites(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertPaginatedResponse(t, PaginatedResponse{
			Offset:     result.Offset,
			Limit:      result.Limit,
			Count:      result.Count,
			TotalCount: result.TotalCount,
		}, PaginatedResponse{
			Offset:     0,
			Limit:      25,
			Count:      1,
			TotalCount: 1,
		})

		if len(result.Data) != 1 {
			t.Fatalf("expected 1 site, got %d", len(result.Data))
		}

		site := result.Data[0]
		if site.ID != expectedSites[0].ID {
			t.Errorf("expected site ID %s, got %s", expectedSites[0].ID, site.ID)
		}
		if site.Name != expectedSites[0].Name {
			t.Errorf("expected site name %s, got %s", expectedSites[0].Name, site.Name)
		}
	})

	t.Run("with pagination parameters", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(200, ListSitesResponse{
			Offset:     50,
			Limit:      10,
			Count:      0,
			TotalCount: 100,
			Data:       []Site{},
		})

		params := &ListSitesParams{
			Offset: 50,
			Limit:  10,
		}

		result, err := client.ListSites(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		assertPaginatedResponse(t, PaginatedResponse{
			Offset:     result.Offset,
			Limit:      result.Limit,
			Count:      result.Count,
			TotalCount: result.TotalCount,
		}, PaginatedResponse{
			Offset:     50,
			Limit:      10,
			Count:      0,
			TotalCount: 100,
		})
	})

	t.Run("invalid limit", func(t *testing.T) {
		client, _ := newTestClient(t, testBaseURL)

		params := &ListSitesParams{
			Limit: 201, // Exceeds maximum of 200
		}

		_, err := client.ListSites(ctx, params)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		const expectedErr = "limit must be between 0 and 200"
		if err.Error() != expectedErr {
			t.Errorf("expected error message %q, got %q", expectedErr, err.Error())
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(401, Error{
			Status:     401,
			StatusName: "Unauthorized",
			Message:    "Invalid credentials",
		})

		_, err := client.ListSites(ctx, nil)
		assertErrorResponse(t, err, 401, "Invalid credentials")
	})
}
