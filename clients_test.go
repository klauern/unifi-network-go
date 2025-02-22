package unifi

import (
	"context"
	"testing"
)

func TestClient_ListNetworkClients(t *testing.T) {
	ctx := context.Background()

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		expectedClients := []NetworkClient{
			{
				ID:             "abc123",
				Name:           "Test Client",
				MACAddress:     "00:11:22:33:44:55",
				IPAddress:      "192.168.1.100",
				Type:           "WIRED",
				UplinkDeviceID: "switch1",
				LastSeen:       1234567890,
				IsWired:        true,
				IsGuest:        false,
				RxBytes:        1000,
				TxBytes:        2000,
				RxRate:         50.5,
				TxRate:         75.2,
			},
		}

		mock.response = mockResponse(200, ListNetworkClientsResponse{
			Offset:     0,
			Limit:      25,
			Count:      1,
			TotalCount: 100,
			Data:       expectedClients,
		})

		result, err := client.ListNetworkClients(ctx, testSiteID, &ListNetworkClientsParams{
			Limit: 25,
		})
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
			TotalCount: 100,
		})

		if len(result.Data) != 1 {
			t.Fatalf("expected 1 client, got %d", len(result.Data))
		}

		resultClient := result.Data[0]
		if resultClient.ID != expectedClients[0].ID {
			t.Errorf("expected client ID %s, got %s", expectedClients[0].ID, resultClient.ID)
		}
		if resultClient.MACAddress != expectedClients[0].MACAddress {
			t.Errorf("expected client MAC %s, got %s", expectedClients[0].MACAddress, resultClient.MACAddress)
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		client, _ := newTestClient(t, testBaseURL)

		_, err := client.ListNetworkClients(ctx, testSiteID, &ListNetworkClientsParams{
			Limit: 201, // Exceeds maximum of 200
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "limit must be between 0 and 200" {
			t.Errorf("expected error message %q, got %q", "limit must be between 0 and 200", err.Error())
		}
	})

	t.Run("site not found", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Site not found",
		})

		_, err := client.ListNetworkClients(ctx, "nonexistent", nil)
		assertErrorResponse(t, err, 404, "Site not found")
	})
}

func TestClient_GetNetworkClient(t *testing.T) {
	ctx := context.Background()
	clientID := "abc123"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		expectedClient := NetworkClient{
			ID:             clientID,
			Name:           "Test Client",
			MACAddress:     "00:11:22:33:44:55",
			IPAddress:      "192.168.1.100",
			Type:           "WIRED",
			UplinkDeviceID: "switch1",
			LastSeen:       1234567890,
			IsWired:        true,
			IsGuest:        false,
			RxBytes:        1000,
			TxBytes:        2000,
			RxRate:         50.5,
			TxRate:         75.2,
		}

		mock.response = mockResponse(200, struct {
			Data []NetworkClient `json:"data"`
		}{
			Data: []NetworkClient{expectedClient},
		})

		result, err := client.GetNetworkClient(ctx, testSiteID, clientID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID != expectedClient.ID {
			t.Errorf("expected client ID %s, got %s", expectedClient.ID, result.ID)
		}
		if result.MACAddress != expectedClient.MACAddress {
			t.Errorf("expected client MAC %s, got %s", expectedClient.MACAddress, result.MACAddress)
		}
	})

	t.Run("client not found", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(200, struct {
			Data []NetworkClient `json:"data"`
		}{
			Data: []NetworkClient{},
		})

		_, err := client.GetNetworkClient(ctx, testSiteID, "nonexistent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "network client not found: nonexistent" {
			t.Errorf("expected error message %q, got %q", "network client not found: nonexistent", err.Error())
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, testBaseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Site not found",
		})

		_, err := client.GetNetworkClient(ctx, "nonexistent", clientID)
		assertErrorResponse(t, err, 404, "Site not found")
	})
}
