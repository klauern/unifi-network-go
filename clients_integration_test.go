package unifi

import (
	"context"
	"testing"
)

func TestIntegration_ListNetworkClients(t *testing.T) {
	client := newIntegrationTestClient(t)
	ctx := context.Background()

	// First get a site ID to use for testing
	sites, err := client.ListSites(ctx, nil)
	if err != nil {
		t.Fatalf("failed to get sites for testing: %v", err)
	}
	if len(sites.Data) == 0 {
		t.Fatal("no sites available for testing")
	}
	siteID := sites.Data[0].ID

	t.Run("list all clients", func(t *testing.T) {
		result, err := client.ListNetworkClients(ctx, siteID, nil)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Basic validation
		if result.Count < 0 {
			t.Error("expected non-negative count")
		}
		if int(result.Count) != len(result.Data) {
			t.Errorf("count mismatch: got %d clients but count is %d", len(result.Data), result.Count)
		}

		// Validate client fields for each client
		for i, client := range result.Data {
			if client.ID == "" {
				t.Errorf("client %d: empty ID", i)
			}
			if client.MACAddress == "" {
				t.Errorf("client %d: empty MAC address", i)
			}

			// Validate connection type
			switch client.Type {
			case "WIRED", "WIRELESS", "VPN", "":
				// Valid types (empty is allowed as it might not be set for all clients)
			default:
				t.Errorf("client %d: invalid type: %s", i, client.Type)
			}

			// Basic field validation
			if client.Name == "" {
				t.Errorf("client %d: empty name", i)
			}
			if client.ConnectedAt == "" {
				t.Errorf("client %d: empty connectedAt", i)
			}
			if client.UplinkDeviceID == "" {
				t.Errorf("client %d: empty uplinkDeviceId", i)
			}
		}
	})

	t.Run("with pagination", func(t *testing.T) {
		// Request first page with small limit
		params := &ListNetworkClientsParams{
			Limit: 1,
		}
		result, err := client.ListNetworkClients(ctx, siteID, params)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		// Validate pagination
		if result.Limit != 1 {
			t.Errorf("expected limit 1, got %d", result.Limit)
		}
		if result.Count > 1 {
			t.Errorf("expected at most 1 client, got %d", result.Count)
		}
		if result.TotalCount < result.Count {
			t.Errorf("total count %d is less than count %d", result.TotalCount, result.Count)
		}

		// If we have more than one client total, test offset
		if result.TotalCount > 1 {
			// Get second page
			params.Offset = 1
			secondPage, err := client.ListNetworkClients(ctx, siteID, params)
			if err != nil {
				t.Fatalf("failed to get second page: %v", err)
			}

			if len(secondPage.Data) > 0 && len(result.Data) > 0 {
				if secondPage.Data[0].MACAddress == result.Data[0].MACAddress {
					t.Error("second page returned same client as first page")
				}
			}
		}
	})

	t.Run("with invalid site ID", func(t *testing.T) {
		_, err := client.ListNetworkClients(ctx, "invalid-site-id", nil)
		if err == nil {
			t.Error("expected error with invalid site ID, got nil")
		}
	})

	t.Run("with invalid limit", func(t *testing.T) {
		params := &ListNetworkClientsParams{
			Limit: 201, // Exceeds maximum of 200
		}
		_, err := client.ListNetworkClients(ctx, siteID, params)
		if err == nil {
			t.Error("expected error with invalid limit, got nil")
		}
	})
}
