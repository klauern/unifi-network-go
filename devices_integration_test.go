package unifi

import (
	"context"
	"fmt"
	"testing"
)

func TestIntegration_ListDevices(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := newIntegrationTestClient(t)

	// First get a site ID to use for testing
	sites, err := client.ListSites(ctx, nil)
	if err != nil {
		t.Fatalf("failed to get sites for testing: %v", err)
	}
	if len(sites.Data) == 0 {
		t.Fatal("no sites available for testing")
	}
	siteID := sites.Data[0].ID

	t.Run("basic listing", func(t *testing.T) {
		resp, err := client.ListDevices(ctx, siteID, nil)
		if err != nil {
			t.Fatalf("ListDevices failed: %v", err)
		}

		// Basic response validation
		if resp == nil {
			t.Fatal("expected non-nil response")
		}
		if resp.Data == nil {
			t.Fatal("expected non-nil Data field")
		}

		// Debug: Print the first device
		if len(resp.Data) > 0 {
			fmt.Printf("First device: %+v\n", resp.Data[0])
		}

		// Validate device fields
		for _, device := range resp.Data {
			// Required fields
			if device.ID == "" {
				t.Error("device ID is required")
			}
			if device.MAC == "" {
				t.Error("device MAC is required")
			}
			if device.Type == "" {
				t.Errorf("device type is required (device: %+v)", device)
			}
			if device.Model == "" {
				t.Error("device model is required")
			}

			// State validation
			if device.State != "ONLINE" && device.State != "OFFLINE" {
				t.Errorf("invalid device state: %s", device.State)
			}

			// Logical field relationships
			if device.State == "ONLINE" {
				if device.IP == "" {
					t.Error("online device should have IP address")
				}
			}
		}
	})

	t.Run("pagination", func(t *testing.T) {
		// First page
		firstPage, err := client.ListDevices(ctx, siteID, &ListDevicesParams{
			Limit: 1,
		})
		if err != nil {
			t.Fatalf("ListDevices failed: %v", err)
		}

		if firstPage.Limit != 1 {
			t.Errorf("expected limit 1, got %d", firstPage.Limit)
		}

		// If we have more than one device, test pagination
		if firstPage.TotalCount > 1 {
			secondPage, err := client.ListDevices(ctx, siteID, &ListDevicesParams{
				Limit:  1,
				Offset: 1,
			})
			if err != nil {
				t.Fatalf("ListDevices failed: %v", err)
			}

			if len(secondPage.Data) == 0 {
				t.Error("expected at least one device in second page")
			}

			if len(firstPage.Data) > 0 && len(secondPage.Data) > 0 {
				if firstPage.Data[0].ID == secondPage.Data[0].ID {
					t.Error("expected different devices in different pages")
				}
			}
		}
	})

	t.Run("device type filter", func(t *testing.T) {
		// Get all devices first to find a valid type
		allDevices, err := client.ListDevices(ctx, siteID, nil)
		if err != nil {
			t.Fatalf("ListDevices failed: %v", err)
		}

		if len(allDevices.Data) > 0 {
			deviceType := allDevices.Data[0].Type
			filtered, err := client.ListDevices(ctx, siteID, &ListDevicesParams{
				Type: deviceType,
			})
			if err != nil {
				t.Fatalf("ListDevices failed: %v", err)
			}

			for _, device := range filtered.Data {
				// Check if the device has the requested type in its features
				hasType := false
				for _, feature := range device.Features {
					if feature == deviceType {
						hasType = true
						break
					}
				}
				if !hasType {
					t.Errorf("expected device to have feature %s, got %v", deviceType, device.Features)
				}
			}
		}
	})

	t.Run("error cases", func(t *testing.T) {
		// Invalid site ID
		_, err := client.ListDevices(ctx, "nonexistent", nil)
		if err == nil {
			t.Error("expected error for invalid site ID")
		}

		// Invalid limit (over 200)
		_, err = client.ListDevices(ctx, siteID, &ListDevicesParams{
			Limit: 201,
		})
		if err == nil {
			t.Error("expected error for limit > 200")
		}

		// Invalid offset (negative)
		_, err = client.ListDevices(ctx, siteID, &ListDevicesParams{
			Offset: -1,
		})
		if err == nil {
			t.Error("expected error for negative offset")
		}
	})
}
