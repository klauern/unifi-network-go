package unifi

import (
	"context"
	"testing"
)

func TestClient_CreateHotspotVoucher(t *testing.T) {
	ctx := context.Background()

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		request := &CreateHotspotVoucherRequest{
			Name:                "Test Voucher",
			TimeLimitMinutes:    1440,
			AuthorizeGuestLimit: 2,
			DataUsageLimitMB:    1024,
			RxRateLimitKbps:     1024,
			TxRateLimitKbps:     512,
			Count:               1,
		}

		expectedVoucher := HotspotVoucher{
			ID:                  "abc123",
			CreatedAt:           "2023-01-01T00:00:00Z",
			Name:                request.Name,
			Code:                "WIFI123",
			AuthorizeGuestLimit: request.AuthorizeGuestLimit,
			TimeLimitMinutes:    request.TimeLimitMinutes,
			DataUsageLimitMB:    request.DataUsageLimitMB,
			RxRateLimitKbps:     request.RxRateLimitKbps,
			TxRateLimitKbps:     request.TxRateLimitKbps,
		}

		mock.response = mockResponse(200, CreateHotspotVoucherResponse{
			Data: []HotspotVoucher{expectedVoucher},
		})

		result, err := client.CreateHotspotVoucher(ctx, testSiteID, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Data) != 1 {
			t.Fatalf("expected 1 voucher, got %d", len(result.Data))
		}

		voucher := result.Data[0]
		if voucher.ID != expectedVoucher.ID {
			t.Errorf("expected voucher ID %s, got %s", expectedVoucher.ID, voucher.ID)
		}
		if voucher.Code != expectedVoucher.Code {
			t.Errorf("expected voucher code %s, got %s", expectedVoucher.Code, voucher.Code)
		}
		if voucher.Name != request.Name {
			t.Errorf("expected voucher name %s, got %s", request.Name, voucher.Name)
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(400, Error{
			Status:     400,
			StatusName: "Bad Request",
			Message:    "Invalid parameters",
		})

		_, err := client.CreateHotspotVoucher(ctx, testSiteID, &CreateHotspotVoucherRequest{})
		assertErrorResponse(t, err, 400, "Invalid parameters")
	})
}

func TestClient_GetHotspotVoucher(t *testing.T) {
	ctx := context.Background()
	voucherID := "abc123"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		expectedVoucher := HotspotVoucher{
			ID:               voucherID,
			CreatedAt:        "2023-01-01T00:00:00Z",
			Name:             "Test Voucher",
			Code:             "WIFI123",
			TimeLimitMinutes: 60,
		}

		mock.response = mockResponse(200, struct {
			Data []HotspotVoucher `json:"data"`
		}{
			Data: []HotspotVoucher{expectedVoucher},
		})

		result, err := client.GetHotspotVoucher(ctx, testSiteID, voucherID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID != expectedVoucher.ID {
			t.Errorf("expected voucher ID %s, got %s", expectedVoucher.ID, result.ID)
		}
		if result.Code != expectedVoucher.Code {
			t.Errorf("expected voucher code %s, got %s", expectedVoucher.Code, result.Code)
		}
	})

	t.Run("voucher not found", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(200, struct {
			Data []HotspotVoucher `json:"data"`
		}{
			Data: []HotspotVoucher{},
		})

		_, err := client.GetHotspotVoucher(ctx, testSiteID, "nonexistent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "voucher not found: nonexistent" {
			t.Errorf("expected error message %q, got %q", "voucher not found: nonexistent", err.Error())
		}
	})
}

func TestClient_DeleteHotspotVoucher(t *testing.T) {
	ctx := context.Background()
	voucherID := "abc123"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)
		mock.response = mockResponse(200, nil)

		err := client.DeleteHotspotVoucher(ctx, testSiteID, voucherID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Voucher not found",
		})

		err := client.DeleteHotspotVoucher(ctx, testSiteID, "nonexistent")
		assertErrorResponse(t, err, 404, "Voucher not found")
	})
}

func TestClient_GenerateHotspotVouchers(t *testing.T) {
	ctx := context.Background()

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		request := &GenerateHotspotVouchersRequest{
			Count:               1,
			Name:                "Test Voucher",
			AuthorizeGuestLimit: 1,
			TimeLimitMinutes:    1,
			DataUsageLimitMB:    1,
			RxRateLimitKbps:     2,
			TxRateLimitKbps:     2,
		}

		mock.response = mockResponse(201, GenerateHotspotVouchersResponse{
			Data: []HotspotVoucher{
				{
					ID:                  "abc123",
					CreatedAt:           "2024-02-22T12:00:00Z",
					Name:                request.Name,
					Code:                "ABC123",
					AuthorizeGuestLimit: request.AuthorizeGuestLimit,
					TimeLimitMinutes:    request.TimeLimitMinutes,
					DataUsageLimitMB:    request.DataUsageLimitMB,
					RxRateLimitKbps:     request.RxRateLimitKbps,
					TxRateLimitKbps:     request.TxRateLimitKbps,
				},
			},
		})

		result, err := client.GenerateHotspotVouchers(ctx, testSiteID, request)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Data) != 1 {
			t.Fatalf("expected 1 voucher, got %d", len(result.Data))
		}

		voucher := result.Data[0]
		if voucher.ID != "abc123" {
			t.Errorf("expected voucher ID abc123, got %s", voucher.ID)
		}
		if voucher.Code != "ABC123" {
			t.Errorf("expected voucher code ABC123, got %s", voucher.Code)
		}
		if voucher.Name != request.Name {
			t.Errorf("expected voucher name %s, got %s", request.Name, voucher.Name)
		}
	})

	t.Run("validation errors", func(t *testing.T) {
		client, _ := newTestClient(t, testBaseURL)

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
				_, err := client.GenerateHotspotVouchers(ctx, testSiteID, tt.request)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.wantErr {
					t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
				}
			})
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(400, Error{
			Status:     400,
			StatusName: "Bad Request",
			Message:    "Invalid parameters",
		})

		request := &GenerateHotspotVouchersRequest{
			Count:            1,
			Name:             "Test",
			TimeLimitMinutes: 1,
		}

		_, err := client.GenerateHotspotVouchers(ctx, testSiteID, request)
		assertErrorResponse(t, err, 400, "Invalid parameters")
	})
}

func TestClient_GetVoucherDetails(t *testing.T) {
	ctx := context.Background()
	voucherID := "4997eeca-0276-4993-bfeb-53cbbbaa4f00"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		expectedVoucher := HotspotVoucher{
			ID:                  voucherID,
			CreatedAt:           "2019-06-24T14:15:22Z",
			Name:                "hotel-guest",
			Code:                "ABC123XYZ",
			AuthorizeGuestLimit: 1,
			AuthorizeGuestCount: 0,
			ActivatedAt:         "2019-06-24T14:15:22Z",
			ExpiresAt:           "2019-06-24T14:15:22Z",
			Expired:             true,
			TimeLimitMinutes:    1440,
			DataUsageLimitMB:    1024,
			RxRateLimitKbps:     1000,
			TxRateLimitKbps:     1000,
		}

		mock.response = mockResponse(200, GetVoucherDetailsResponse{
			Data: []HotspotVoucher{expectedVoucher},
		})

		result, err := client.GetVoucherDetails(ctx, testSiteID, voucherID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID != expectedVoucher.ID {
			t.Errorf("expected voucher ID %s, got %s", expectedVoucher.ID, result.ID)
		}
		if result.Code != expectedVoucher.Code {
			t.Errorf("expected voucher code %s, got %s", expectedVoucher.Code, result.Code)
		}
		if result.Name != expectedVoucher.Name {
			t.Errorf("expected voucher name %s, got %s", expectedVoucher.Name, result.Name)
		}
		if result.CreatedAt != expectedVoucher.CreatedAt {
			t.Errorf("expected voucher createdAt %s, got %s", expectedVoucher.CreatedAt, result.CreatedAt)
		}
		if result.TimeLimitMinutes != expectedVoucher.TimeLimitMinutes {
			t.Errorf("expected voucher timeLimitMinutes %d, got %d", expectedVoucher.TimeLimitMinutes, result.TimeLimitMinutes)
		}
	})

	t.Run("validation errors", func(t *testing.T) {
		client, _ := newTestClient(t, testBaseURL)

		tests := []struct {
			name      string
			siteID    string
			voucherID string
			wantErr   string
		}{
			{
				name:      "missing site ID",
				siteID:    "",
				voucherID: "test",
				wantErr:   "siteId is required",
			},
			{
				name:      "missing voucher ID",
				siteID:    "test",
				voucherID: "",
				wantErr:   "voucherId is required",
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				_, err := client.GetVoucherDetails(ctx, tt.siteID, tt.voucherID)
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				if err.Error() != tt.wantErr {
					t.Errorf("expected error %q, got %q", tt.wantErr, err.Error())
				}
			})
		}
	})

	t.Run("voucher not found", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(200, GetVoucherDetailsResponse{
			Data: []HotspotVoucher{},
		})

		_, err := client.GetVoucherDetails(ctx, testSiteID, "nonexistent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "voucher not found: nonexistent" {
			t.Errorf("expected error message %q, got %q", "voucher not found: nonexistent", err.Error())
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Voucher not found",
		})

		_, err := client.GetVoucherDetails(ctx, testSiteID, "nonexistent")
		assertErrorResponse(t, err, 404, "Voucher not found")
	})
}
