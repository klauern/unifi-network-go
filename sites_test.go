package unifi

import (
	"context"
	"errors"
	"testing"
)

func TestClient_ListSites(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		expectedSites := []Site{
			{
				ID:          "default",
				Name:        "Default",
				Description: "Default site",
				Role:        "admin",
				Hidden:      false,
				NoDelete:    true,
			},
		}

		mock.response = mockResponse(200, ListSitesResponse{
			PaginatedResponse: PaginatedResponse{
				Offset:     0,
				Limit:      25,
				Count:      1,
				TotalCount: 1,
			},
			Data: expectedSites,
		})

		result, err := client.ListSites(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

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
		if site.Description != expectedSites[0].Description {
			t.Errorf("expected site description %s, got %s", expectedSites[0].Description, site.Description)
		}
		if site.Role != expectedSites[0].Role {
			t.Errorf("expected site role %s, got %s", expectedSites[0].Role, site.Role)
		}
		if site.Hidden != expectedSites[0].Hidden {
			t.Errorf("expected site hidden %v, got %v", expectedSites[0].Hidden, site.Hidden)
		}
		if site.NoDelete != expectedSites[0].NoDelete {
			t.Errorf("expected site no_delete %v, got %v", expectedSites[0].NoDelete, site.NoDelete)
		}
	})

	t.Run("with pagination parameters", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(200, ListSitesResponse{
			PaginatedResponse: PaginatedResponse{
				Offset:     50,
				Limit:      10,
				Count:      0,
				TotalCount: 100,
			},
			Data: []Site{},
		})

		params := &ListSitesParams{
			Offset: 50,
			Limit:  10,
		}

		result, err := client.ListSites(ctx, params)
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

	t.Run("invalid limit", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		params := &ListSitesParams{
			Limit: 201, // Exceeds maximum of 200
		}

		_, err := client.ListSites(ctx, params)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "limit must be between 0 and 200" {
			t.Errorf("expected error message %q, got %q", "limit must be between 0 and 200", err.Error())
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(401, Error{
			Status:     401,
			StatusName: "Unauthorized",
			Message:    "Invalid credentials",
		})

		_, err := client.ListSites(ctx, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("expected *Error, got %T", err)
		}
	})
}

func TestClient_GetSite(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()
	siteID := "default"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		expectedSite := Site{
			ID:          siteID,
			Name:        "Default",
			Description: "Default site",
			Role:        "admin",
			Hidden:      false,
			NoDelete:    true,
		}

		mock.response = mockResponse(200, struct {
			Data []Site `json:"data"`
		}{
			Data: []Site{expectedSite},
		})

		result, err := client.GetSite(ctx, siteID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID != expectedSite.ID {
			t.Errorf("expected site ID %s, got %s", expectedSite.ID, result.ID)
		}
		if result.Name != expectedSite.Name {
			t.Errorf("expected site name %s, got %s", expectedSite.Name, result.Name)
		}
		if result.Description != expectedSite.Description {
			t.Errorf("expected site description %s, got %s", expectedSite.Description, result.Description)
		}
		if result.Role != expectedSite.Role {
			t.Errorf("expected site role %s, got %s", expectedSite.Role, result.Role)
		}
		if result.Hidden != expectedSite.Hidden {
			t.Errorf("expected site hidden %v, got %v", expectedSite.Hidden, result.Hidden)
		}
		if result.NoDelete != expectedSite.NoDelete {
			t.Errorf("expected site no_delete %v, got %v", expectedSite.NoDelete, result.NoDelete)
		}
	})

	t.Run("empty site ID", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		_, err := client.GetSite(ctx, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "siteID cannot be empty" {
			t.Errorf("expected error message %q, got %q", "siteID cannot be empty", err.Error())
		}
	})

	t.Run("site not found", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(200, struct {
			Data []Site `json:"data"`
		}{
			Data: []Site{},
		})

		_, err := client.GetSite(ctx, "nonexistent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "site not found: nonexistent" {
			t.Errorf("expected error message %q, got %q", "site not found: nonexistent", err.Error())
		}
	})
}
