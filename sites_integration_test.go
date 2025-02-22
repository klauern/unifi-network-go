package unifi

import (
	"context"
	"testing"
)

func TestIntegration_ListSites(t *testing.T) {
	client := newIntegrationTestClient(t)
	ctx := context.Background()

	t.Run("list all sites", func(t *testing.T) {
		result, err := client.ListSites(ctx, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Basic validation
		if result.Count == 0 {
			t.Error("expected at least one site")
		}
		if int(result.Count) != len(result.Data) {
			t.Errorf("count mismatch: got %d sites but count is %d", len(result.Data), result.Count)
		}

		// Validate site fields
		for i, site := range result.Data {
			if site.ID == "" {
				t.Errorf("site %d: empty ID", i)
			}
			if site.Name == "" {
				t.Errorf("site %d: empty Name", i)
			}
		}
	})

	t.Run("with pagination", func(t *testing.T) {
		// Request first page with small limit
		params := &ListSitesParams{
			Limit: 1,
		}
		result, err := client.ListSites(ctx, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Validate pagination
		if result.Limit != 1 {
			t.Errorf("expected limit 1, got %d", result.Limit)
		}
		if result.Count > 1 {
			t.Errorf("expected at most 1 site, got %d", result.Count)
		}
		if result.TotalCount < result.Count {
			t.Errorf("total count %d is less than count %d", result.TotalCount, result.Count)
		}
	})
}
