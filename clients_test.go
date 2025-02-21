package unifi

import (
	"context"
	"errors"
	"testing"
)

func TestClient_ListNetworkClients(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()
	siteID := "default"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		expectedClients := []NetworkClient{
			{
				ID:             "abc123",
				Name:           "test-device",
				ConnectedAt:    "2024-01-01T00:00:00Z",
				IPAddress:      "192.168.1.100",
				Type:           "WIRED",
				MACAddress:     "00:11:22:33:44:55",
				UplinkDeviceID: "switch1",
				SiteID:         siteID,
				Network:        "LAN",
				NetworkName:    "Default",
				IsWired:        true,
				IsGuest:        false,
				DeviceID:       "uap1",
				DeviceName:     "Test AP",
				DeviceMAC:      "aa:bb:cc:dd:ee:ff",
				RxBytes:        1000,
				TxBytes:        2000,
				SignalStrength: -65,
				Channel:        36,
				RadioProtocol:  "ac",
				RadioBand:      "5g",
				SSID:           "Test-SSID",
				Blocked:        false,
				Authorized:     true,
			},
		}

		mock.response = mockResponse(200, ListNetworkClientsResponse{
			PaginatedResponse: PaginatedResponse{
				Offset:     0,
				Limit:      25,
				Count:      1,
				TotalCount: 1,
			},
			Data: expectedClients,
		})

		result, err := client.ListNetworkClients(ctx, siteID, &ListNetworkClientsParams{
			Limit:       100,
			Type:        "wired",
			WithinHours: 24,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Data) != 1 {
			t.Fatalf("expected 1 client, got %d", len(result.Data))
		}

		networkClient := result.Data[0]
		if networkClient.ID != expectedClients[0].ID {
			t.Errorf("expected client ID %s, got %s", expectedClients[0].ID, networkClient.ID)
		}
		if networkClient.Name != expectedClients[0].Name {
			t.Errorf("expected client name %s, got %s", expectedClients[0].Name, networkClient.Name)
		}
		if networkClient.Type != expectedClients[0].Type {
			t.Errorf("expected client type %s, got %s", expectedClients[0].Type, networkClient.Type)
		}
		if networkClient.MACAddress != expectedClients[0].MACAddress {
			t.Errorf("expected client MAC %s, got %s", expectedClients[0].MACAddress, networkClient.MACAddress)
		}
		if networkClient.IPAddress != expectedClients[0].IPAddress {
			t.Errorf("expected client IP %s, got %s", expectedClients[0].IPAddress, networkClient.IPAddress)
		}
	})

	t.Run("empty site ID", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		_, err := client.ListNetworkClients(ctx, "", nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "siteId is required" {
			t.Errorf("expected error message %q, got %q", "siteId is required", err.Error())
		}
	})

	t.Run("invalid limit", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		_, err := client.ListNetworkClients(ctx, siteID, &ListNetworkClientsParams{
			Limit: 201, // Exceeds maximum of 200
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "limit must be between 0 and 200" {
			t.Errorf("expected error message %q, got %q", "limit must be between 0 and 200", err.Error())
		}
	})

	t.Run("invalid type", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		_, err := client.ListNetworkClients(ctx, siteID, &ListNetworkClientsParams{
			Type: "invalid",
		})
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "type must be one of: all, wired, wireless" {
			t.Errorf("expected error message %q, got %q", "type must be one of: all, wired, wireless", err.Error())
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Site not found",
		})

		_, err := client.ListNetworkClients(ctx, "nonexistent", nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("expected *Error, got %T", err)
		}
	})
}

func TestClient_GetNetworkClient(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()
	siteID := "default"
	clientID := "abc123"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		expectedClient := NetworkClient{
			ID:             clientID,
			Name:           "test-device",
			ConnectedAt:    "2024-01-01T00:00:00Z",
			IPAddress:      "192.168.1.100",
			Type:           "WIRED",
			MACAddress:     "00:11:22:33:44:55",
			UplinkDeviceID: "switch1",
			SiteID:         siteID,
			Network:        "LAN",
			NetworkName:    "Default",
			IsWired:        true,
			IsGuest:        false,
			DeviceID:       "uap1",
			DeviceName:     "Test AP",
			DeviceMAC:      "aa:bb:cc:dd:ee:ff",
			SignalStrength: -65,
			Channel:        36,
			RadioProtocol:  "ac",
			RadioBand:      "5g",
			SSID:           "Test-SSID",
			Blocked:        false,
			Authorized:     true,
		}

		mock.response = mockResponse(200, struct {
			Data []NetworkClient `json:"data"`
		}{
			Data: []NetworkClient{expectedClient},
		})

		result, err := client.GetNetworkClient(ctx, siteID, clientID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID != expectedClient.ID {
			t.Errorf("expected client ID %s, got %s", expectedClient.ID, result.ID)
		}
		if result.Name != expectedClient.Name {
			t.Errorf("expected client name %s, got %s", expectedClient.Name, result.Name)
		}
		if result.Type != expectedClient.Type {
			t.Errorf("expected client type %s, got %s", expectedClient.Type, result.Type)
		}
		if result.MACAddress != expectedClient.MACAddress {
			t.Errorf("expected client MAC %s, got %s", expectedClient.MACAddress, result.MACAddress)
		}
	})

	t.Run("empty site ID", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		_, err := client.GetNetworkClient(ctx, "", clientID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "siteId is required" {
			t.Errorf("expected error message %q, got %q", "siteId is required", err.Error())
		}
	})

	t.Run("empty client ID", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		_, err := client.GetNetworkClient(ctx, siteID, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "clientId is required" {
			t.Errorf("expected error message %q, got %q", "clientId is required", err.Error())
		}
	})

	t.Run("client not found", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(200, struct {
			Data []NetworkClient `json:"data"`
		}{
			Data: []NetworkClient{},
		})

		_, err := client.GetNetworkClient(ctx, siteID, "nonexistent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "network client not found: nonexistent" {
			t.Errorf("expected error message %q, got %q", "network client not found: nonexistent", err.Error())
		}
	})
}

func TestClient_BlockUnblockNetworkClient(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()
	siteID := "default"
	clientID := "abc123"

	t.Run("block client", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)
		mock.response = mockResponse(200, nil)

		err := client.BlockNetworkClient(ctx, siteID, clientID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("block client empty site ID", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		err := client.BlockNetworkClient(ctx, "", clientID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "siteId is required" {
			t.Errorf("expected error message %q, got %q", "siteId is required", err.Error())
		}
	})

	t.Run("block client empty client ID", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		err := client.BlockNetworkClient(ctx, siteID, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "clientId is required" {
			t.Errorf("expected error message %q, got %q", "clientId is required", err.Error())
		}
	})

	t.Run("unblock client", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)
		mock.response = mockResponse(200, nil)

		err := client.UnblockNetworkClient(ctx, siteID, clientID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("unblock client empty site ID", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		err := client.UnblockNetworkClient(ctx, "", clientID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "siteId is required" {
			t.Errorf("expected error message %q, got %q", "siteId is required", err.Error())
		}
	})

	t.Run("unblock client empty client ID", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		err := client.UnblockNetworkClient(ctx, siteID, "")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "clientId is required" {
			t.Errorf("expected error message %q, got %q", "clientId is required", err.Error())
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Client not found",
		})

		err := client.BlockNetworkClient(ctx, siteID, "nonexistent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("expected *Error, got %T", err)
		}
	})
}
