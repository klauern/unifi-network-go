package unifi

import (
	"context"
	"testing"
)

func TestIntegration_HotspotVouchers(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	ctx := context.Background()
	client := newIntegrationTestClient(t)

	// Get a valid site ID
	sites, err := client.ListSites(ctx, nil)
	if err != nil {
		t.Fatalf("failed to list sites: %v", err)
	}
	if len(sites.Data) == 0 {
		t.Fatal("no sites found")
	}
	siteID := sites.Data[0].ID

	t.Run("list_vouchers", func(t *testing.T) {
		t.Run("basic_listing", func(t *testing.T) {
			resp, err := client.ListHotspotVouchers(ctx, siteID, nil)
			if err != nil {
				t.Fatalf("failed to list vouchers: %v", err)
			}
			if resp == nil {
				t.Fatal("response is nil")
			}
			if resp.Data == nil {
				t.Fatal("response data is nil")
			}
			t.Logf("Found %d vouchers", len(resp.Data))
		})

		t.Run("pagination", func(t *testing.T) {
			params := &ListHotspotVouchersParams{
				Limit:  1,
				Offset: 0,
			}
			resp, err := client.ListHotspotVouchers(ctx, siteID, params)
			if err != nil {
				t.Fatalf("failed to list vouchers: %v", err)
			}
			if resp == nil {
				t.Fatal("response is nil")
			}
			if resp.Data == nil {
				t.Fatal("response data is nil")
			}
		})
	})

	// Create a test voucher
	createReq := &GenerateHotspotVouchersRequest{
		Count:               1,
		Name:                "Integration Test Voucher",
		TimeLimitMinutes:    60,
		AuthorizeGuestLimit: 1,
		DataUsageLimitMB:    1024,
		RxRateLimitKbps:     2,
		TxRateLimitKbps:     2,
	}

	t.Run("generate_voucher", func(t *testing.T) {
		resp, err := client.GenerateHotspotVouchers(ctx, siteID, createReq)
		if err != nil {
			t.Fatalf("failed to generate voucher: %v", err)
		}
		if resp == nil {
			t.Fatal("response is nil")
		}
		if len(resp.Data) == 0 {
			t.Fatal("no vouchers generated")
		}

		voucher := resp.Data[0]
		t.Logf("Generated voucher with ID: %s", voucher.ID)

		// Test getting voucher details
		t.Run("get_voucher_details", func(t *testing.T) {
			details, err := client.GetVoucherDetails(ctx, siteID, voucher.ID)
			if err != nil {
				t.Fatalf("failed to get voucher details: %v", err)
			}
			if details == nil {
				t.Fatal("voucher details is nil")
			}
			if details.ID != voucher.ID {
				t.Errorf("expected voucher ID %s, got %s", voucher.ID, details.ID)
			}
			if details.Name != createReq.Name {
				t.Errorf("expected voucher name %s, got %s", createReq.Name, details.Name)
			}
			if details.TimeLimitMinutes != createReq.TimeLimitMinutes {
				t.Errorf("expected time limit %d, got %d", createReq.TimeLimitMinutes, details.TimeLimitMinutes)
			}
			if details.AuthorizeGuestLimit != createReq.AuthorizeGuestLimit {
				t.Errorf("expected guest limit %d, got %d", createReq.AuthorizeGuestLimit, details.AuthorizeGuestLimit)
			}
			if details.DataUsageLimitMB != createReq.DataUsageLimitMB {
				t.Errorf("expected data limit %d, got %d", createReq.DataUsageLimitMB, details.DataUsageLimitMB)
			}
			if details.RxRateLimitKbps != createReq.RxRateLimitKbps {
				t.Errorf("expected rx rate limit %d, got %d", createReq.RxRateLimitKbps, details.RxRateLimitKbps)
			}
			if details.TxRateLimitKbps != createReq.TxRateLimitKbps {
				t.Errorf("expected tx rate limit %d, got %d", createReq.TxRateLimitKbps, details.TxRateLimitKbps)
			}
		})

		// Clean up by deleting the test voucher
		t.Run("delete_voucher", func(t *testing.T) {
			err := client.DeleteHotspotVoucher(ctx, siteID, voucher.ID)
			if err != nil {
				t.Fatalf("failed to delete voucher: %v", err)
			}

			// Verify deletion
			_, err = client.GetVoucherDetails(ctx, siteID, voucher.ID)
			if err == nil {
				t.Error("expected error getting deleted voucher, got nil")
			}
		})
	})

	t.Run("error_cases", func(t *testing.T) {
		t.Run("invalid_site_id", func(t *testing.T) {
			_, err := client.ListHotspotVouchers(ctx, "invalid-site", nil)
			if err == nil {
				t.Error("expected error for invalid site ID, got nil")
			}
		})

		t.Run("invalid_voucher_id", func(t *testing.T) {
			_, err := client.GetVoucherDetails(ctx, siteID, "invalid-voucher")
			if err == nil {
				t.Error("expected error for invalid voucher ID, got nil")
			}
		})

		t.Run("invalid_request", func(t *testing.T) {
			tests := []struct {
				name    string
				request *GenerateHotspotVouchersRequest
				wantErr string
			}{
				{
					name:    "nil request",
					request: nil,
					wantErr: "request cannot be nil",
				},
				{
					name:    "missing name",
					request: &GenerateHotspotVouchersRequest{Count: 1, TimeLimitMinutes: 1},
					wantErr: "name is required",
				},
				{
					name:    "count too low",
					request: &GenerateHotspotVouchersRequest{Count: 0, Name: "Test", TimeLimitMinutes: 1},
					wantErr: "count must be between 1 and 10000",
				},
				{
					name:    "count too high",
					request: &GenerateHotspotVouchersRequest{Count: 10001, Name: "Test", TimeLimitMinutes: 1},
					wantErr: "count must be between 1 and 10000",
				},
				{
					name:    "time limit too low",
					request: &GenerateHotspotVouchersRequest{Count: 1, Name: "Test", TimeLimitMinutes: 0},
					wantErr: "timeLimitMinutes must be between 1 and 1000000",
				},
				{
					name:    "time limit too high",
					request: &GenerateHotspotVouchersRequest{Count: 1, Name: "Test", TimeLimitMinutes: 1000001},
					wantErr: "timeLimitMinutes must be between 1 and 1000000",
				},
				{
					name:    "data usage limit too high",
					request: &GenerateHotspotVouchersRequest{Count: 1, Name: "Test", TimeLimitMinutes: 1, DataUsageLimitMB: 1046577},
					wantErr: "dataUsageLimitMBytes must be between 1 and 1046576",
				},
				{
					name:    "rx rate limit too low",
					request: &GenerateHotspotVouchersRequest{Count: 1, Name: "Test", TimeLimitMinutes: 1, RxRateLimitKbps: 1},
					wantErr: "rxRateLimitKbps must be between 2 and 100000",
				},
				{
					name:    "tx rate limit too high",
					request: &GenerateHotspotVouchersRequest{Count: 1, Name: "Test", TimeLimitMinutes: 1, TxRateLimitKbps: 100001},
					wantErr: "txRateLimitKbps must be between 2 and 100000",
				},
			}

			for _, tt := range tests {
				t.Run(tt.name, func(t *testing.T) {
					_, err := client.GenerateHotspotVouchers(ctx, siteID, tt.request)
					if err == nil {
						t.Fatal("expected error, got nil")
					}
					if err.Error() != tt.wantErr {
						t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
					}
				})
			}
		})
	})
}

func validateVoucher(t *testing.T, v HotspotVoucher) {
	t.Helper()

	if v.ID == "" {
		t.Error("voucher ID is empty")
	}
	if v.CreatedAt == "" {
		t.Error("voucher CreatedAt is empty")
	}
	if v.Name == "" {
		t.Error("voucher Name is empty")
	}
	if v.Code == "" {
		t.Error("voucher Code is empty")
	}
	if v.TimeLimitMinutes <= 0 {
		t.Error("voucher TimeLimitMinutes is not positive")
	}
}
