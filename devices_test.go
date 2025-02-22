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

		expectedResponse := ListDevicesResponse{
			PaginatedResponse: PaginatedResponse{
				Count:      1,
				TotalCount: 100,
			},
			Data: []Device{
				{
					ID:         "abc123",
					Name:       "test-device",
					Type:       "uap",
					Model:      "U6-Pro",
					Version:    "6.0.15",
					State:      1,
					IP:         "192.168.1.100",
					MAC:        "00:11:22:33:44:55",
					Disabled:   false,
					SiteID:     siteID,
					Adopted:    true,
					LastSeen:   1234567890,
					Upgradable: false,
				},
			},
		}

		mock.response = mockResponse(200, expectedResponse)

		result, err := client.ListDevices(ctx, siteID, &ListDevicesParams{
			Limit: 25,
		})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if result.Count != expectedResponse.Count {
			t.Errorf("expected count %d, got %d", expectedResponse.Count, result.Count)
		}

		if len(result.Data) != 1 {
			t.Fatalf("expected 1 device, got %d", len(result.Data))
		}

		device := result.Data[0]
		if device.ID != expectedResponse.Data[0].ID {
			t.Errorf("expected device ID %s, got %s", expectedResponse.Data[0].ID, device.ID)
		}
		if device.Name != expectedResponse.Data[0].Name {
			t.Errorf("expected device name %s, got %s", expectedResponse.Data[0].Name, device.Name)
		}
		if device.Type != expectedResponse.Data[0].Type {
			t.Errorf("expected device type %s, got %s", expectedResponse.Data[0].Type, device.Type)
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
			Name:       "test-device",
			Type:       "uap",
			Model:      "U6-Pro",
			Version:    "6.0.15",
			State:      1,
			IP:         "192.168.1.100",
			MAC:        "00:11:22:33:44:55",
			Disabled:   false,
			SiteID:     siteID,
			Adopted:    true,
			LastSeen:   1234567890,
			Upgradable: false,
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
		if result.Name != expectedDevice.Name {
			t.Errorf("expected device name %s, got %s", expectedDevice.Name, result.Name)
		}
		if result.Type != expectedDevice.Type {
			t.Errorf("expected device type %s, got %s", expectedDevice.Type, result.Type)
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
		if err.Error() != "device not found: nonexistent" {
			t.Errorf("expected error message %q, got %q", "device not found: nonexistent", err.Error())
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
