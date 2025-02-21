package unifi

import (
	"context"
	"errors"
	"testing"
)

func TestClient_ListDevices(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()
	siteID := "default"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		expectedDevices := []Device{
			{
				ID:         "abc123",
				MAC:        "00:11:22:33:44:55",
				Model:      "U6-Lite",
				Type:       "uap",
				Name:       "Test AP",
				SiteID:     siteID,
				IP:         "192.168.1.10",
				Version:    "6.0.0",
				Adopted:    true,
				Disabled:   false,
				LastSeen:   1234567890,
				Upgradable: false,
				State:      1,
			},
		}

		mock.response = mockResponse(200, ListDevicesResponse{
			PaginatedResponse: PaginatedResponse{
				Count:      1,
				TotalCount: 1,
			},
			Data: expectedDevices,
		})

		result, err := client.ListDevices(ctx, siteID, &ListDevicesParams{
			Limit: 100,
			Type:  "uap",
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(result.Data) != 1 {
			t.Fatalf("expected 1 device, got %d", len(result.Data))
		}

		device := result.Data[0]
		if device.ID != expectedDevices[0].ID {
			t.Errorf("expected device ID %s, got %s", expectedDevices[0].ID, device.ID)
		}
		if device.MAC != expectedDevices[0].MAC {
			t.Errorf("expected device MAC %s, got %s", expectedDevices[0].MAC, device.MAC)
		}
		if device.Type != expectedDevices[0].Type {
			t.Errorf("expected device type %s, got %s", expectedDevices[0].Type, device.Type)
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Site not found",
		})

		_, err := client.ListDevices(ctx, "nonexistent", nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("expected *Error, got %T", err)
		}
	})
}

func TestClient_GetDevice(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()
	siteID := "default"
	deviceID := "abc123"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		expectedDevice := Device{
			ID:         deviceID,
			MAC:        "00:11:22:33:44:55",
			Model:      "U6-Lite",
			Type:       "uap",
			Name:       "Test AP",
			SiteID:     siteID,
			IP:         "192.168.1.10",
			Version:    "6.0.0",
			Adopted:    true,
			Disabled:   false,
			LastSeen:   1234567890,
			Upgradable: false,
			State:      1,
		}

		mock.response = mockResponse(200, struct {
			Data []Device `json:"data"`
		}{
			Data: []Device{expectedDevice},
		})

		result, err := client.GetDevice(ctx, siteID, deviceID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID != expectedDevice.ID {
			t.Errorf("expected device ID %s, got %s", expectedDevice.ID, result.ID)
		}
		if result.MAC != expectedDevice.MAC {
			t.Errorf("expected device MAC %s, got %s", expectedDevice.MAC, result.MAC)
		}
	})

	t.Run("device not found", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(200, struct {
			Data []Device `json:"data"`
		}{
			Data: []Device{},
		})

		_, err := client.GetDevice(ctx, siteID, "nonexistent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})
}

func TestClient_ExecutePortAction(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()
	siteID := "default"
	deviceID := "abc123"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)
		mock.response = mockResponse(200, nil)

		action := &DevicePortAction{
			PortIDX: 1,
			PortID:  "port1",
			Action:  "reset",
		}

		err := client.ExecutePortAction(ctx, siteID, deviceID, action)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("nil action", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		err := client.ExecutePortAction(ctx, siteID, deviceID, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "action cannot be nil" {
			t.Errorf("expected error message %q, got %q", "action cannot be nil", err.Error())
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Device not found",
		})

		action := &DevicePortAction{
			PortIDX: 1,
			PortID:  "port1",
			Action:  "reset",
		}

		err := client.ExecutePortAction(ctx, siteID, "nonexistent", action)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("expected *Error, got %T", err)
		}
	})
}

func TestClient_ExecuteDeviceAction(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()
	siteID := "default"
	deviceID := "abc123"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)
		mock.response = mockResponse(200, nil)

		action := &DeviceAction{
			Action: "restart",
		}

		err := client.ExecuteDeviceAction(ctx, siteID, deviceID, action)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})

	t.Run("nil action", func(t *testing.T) {
		client, _ := newTestClient(t, baseURL)

		err := client.ExecuteDeviceAction(ctx, siteID, deviceID, nil)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "action cannot be nil" {
			t.Errorf("expected error message %q, got %q", "action cannot be nil", err.Error())
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Device not found",
		})

		action := &DeviceAction{
			Action: "restart",
		}

		err := client.ExecuteDeviceAction(ctx, siteID, "nonexistent", action)
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("expected *Error, got %T", err)
		}
	})
}

func TestClient_GetDeviceStatistics(t *testing.T) {
	baseURL := "https://192.168.1.1"
	ctx := context.Background()
	siteID := "default"
	deviceID := "abc123"

	t.Run("successful request", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		expectedStats := DeviceStatistics{
			ID:      deviceID,
			MAC:     "00:11:22:33:44:55",
			RxBytes: 1000,
			TxBytes: 2000,
			RxRate:  50.5,
			TxRate:  75.2,
			CPU:     25.5,
			Memory:  45.2,
			SystemStats: struct {
				Temperature float64 `json:"temperature"`
				FanLevel    int     `json:"fan_level"`
			}{
				Temperature: 45.5,
				FanLevel:    2,
			},
			Uptime: 3600,
		}

		mock.response = mockResponse(200, struct {
			Data []DeviceStatistics `json:"data"`
		}{
			Data: []DeviceStatistics{expectedStats},
		})

		result, err := client.GetDeviceStatistics(ctx, siteID, deviceID)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.ID != expectedStats.ID {
			t.Errorf("expected stats ID %s, got %s", expectedStats.ID, result.ID)
		}
		if result.MAC != expectedStats.MAC {
			t.Errorf("expected stats MAC %s, got %s", expectedStats.MAC, result.MAC)
		}
		if result.CPU != expectedStats.CPU {
			t.Errorf("expected CPU usage %.2f, got %.2f", expectedStats.CPU, result.CPU)
		}
		if result.Memory != expectedStats.Memory {
			t.Errorf("expected memory usage %.2f, got %.2f", expectedStats.Memory, result.Memory)
		}
		if result.SystemStats.Temperature != expectedStats.SystemStats.Temperature {
			t.Errorf("expected temperature %.2f, got %.2f", expectedStats.SystemStats.Temperature, result.SystemStats.Temperature)
		}
	})

	t.Run("no statistics found", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(200, struct {
			Data []DeviceStatistics `json:"data"`
		}{
			Data: []DeviceStatistics{},
		})

		_, err := client.GetDeviceStatistics(ctx, siteID, deviceID)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
		if err.Error() != "no statistics found for device: "+deviceID {
			t.Errorf("expected error message %q, got %q", "no statistics found for device: "+deviceID, err.Error())
		}
	})

	t.Run("error response", func(t *testing.T) {
		client, mock := newTestClient(t, baseURL)

		mock.response = mockResponse(404, Error{
			Status:     404,
			StatusName: "Not Found",
			Message:    "Device not found",
		})

		_, err := client.GetDeviceStatistics(ctx, siteID, "nonexistent")
		if err == nil {
			t.Fatal("expected error, got nil")
		}

		var apiErr *Error
		if !errors.As(err, &apiErr) {
			t.Errorf("expected *Error, got %T", err)
		}
	})
}
