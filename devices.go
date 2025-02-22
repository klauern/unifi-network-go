package unifi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// Device represents a UniFi network device
type Device struct {
	ID         string   `json:"id"`
	MAC        string   `json:"macAddress"`
	Model      string   `json:"model"`
	Type       string   `json:"-"` // Derived from features
	Features   []string `json:"features"`
	Name       string   `json:"name"`
	SiteID     string   `json:"siteId"`
	IP         string   `json:"ipAddress"`
	Version    string   `json:"version"`
	Adopted    bool     `json:"adopted"`
	Disabled   bool     `json:"disabled"`
	Uptime     int64    `json:"uptime"`
	LastSeen   int64    `json:"lastSeen"`
	Upgradable bool     `json:"upgradable"`
	State      string   `json:"state"`
	LastUplink string   `json:"lastUplink"`
	UplinkMAC  string   `json:"uplinkMac"`
}

// UnmarshalJSON implements json.Unmarshaler to derive Type from Features
func (d *Device) UnmarshalJSON(data []byte) error {
	type Alias Device
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(d),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	// Derive Type from Features
	if len(d.Features) > 0 {
		d.Type = d.Features[0]
	}

	return nil
}

// DevicePortAction represents the action to perform on a device port
type DevicePortAction struct {
	PortIDX int    `json:"portIdx"` // Port index number
	PortID  string `json:"portId"`  // Port identifier
	Action  string `json:"action"`  // Action to perform (e.g., "reset", "enable", "disable")
}

// DeviceAction represents the action to perform on a device
type DeviceAction struct {
	Action string `json:"cmd"` // Action to perform (e.g., "restart", "adopt", "forget")
}

// DeviceStatistics represents the latest statistics for a device
type DeviceStatistics struct {
	ID          string  `json:"_id"`
	MAC         string  `json:"mac"`
	RxBytes     int64   `json:"rx_bytes"`
	TxBytes     int64   `json:"tx_bytes"`
	RxRate      float64 `json:"rx_rate"`
	TxRate      float64 `json:"tx_rate"`
	RxPackets   int64   `json:"rx_packets"`
	TxPackets   int64   `json:"tx_packets"`
	RxErrors    int64   `json:"rx_errors"`
	TxErrors    int64   `json:"tx_errors"`
	RxDropped   int64   `json:"rx_dropped"`
	TxDropped   int64   `json:"tx_dropped"`
	RxMulticast int64   `json:"rx_multicast"`
	TxMulticast int64   `json:"tx_multicast"`
	RxBroadcast int64   `json:"rx_broadcast"`
	TxBroadcast int64   `json:"tx_broadcast"`
	BytesR      int64   `json:"bytes-r"`    // Total bytes in last interval
	RxBytesR    int64   `json:"rx_bytes-r"` // Received bytes in last interval
	TxBytesR    int64   `json:"tx_bytes-r"` // Transmitted bytes in last interval
	CPU         float64 `json:"cpu"`        // CPU usage percentage
	Memory      float64 `json:"mem"`        // Memory usage percentage
	SystemStats struct {
		Temperature float64 `json:"temperature"` // Device temperature
		FanLevel    int     `json:"fan_level"`   // Fan level (if applicable)
	} `json:"system-stats"`
	Uptime int64 `json:"uptime"` // Device uptime in seconds
}

// ListDevicesParams contains parameters for listing devices
type ListDevicesParams struct {
	Offset int    `json:"offset,omitempty"`
	Limit  int    `json:"limit,omitempty"`
	Type   string `json:"type,omitempty"`
}

// ListDevicesResponse represents the response from listing devices
type ListDevicesResponse struct {
	PaginatedResponse
	Data []Device `json:"data"`
}

// ListDevices retrieves a paginated list of devices for a site
func (c *Client) ListDevices(ctx context.Context, siteID string, params *ListDevicesParams) (*ListDevicesResponse, error) {
	const maxLimit = 200

	if params != nil {
		if params.Limit > maxLimit {
			return nil, fmt.Errorf("limit must be between 0 and %d", maxLimit)
		}
		if params.Offset < 0 {
			return nil, fmt.Errorf("offset must be non-negative")
		}
	}

	urlPath := fmt.Sprintf("/api/v1/sites/%s/devices", siteID)

	if params != nil {
		query := url.Values{}
		if params.Offset > 0 {
			query.Set("offset", fmt.Sprint(params.Offset))
		}
		if params.Limit > 0 {
			query.Set("limit", fmt.Sprint(params.Limit))
		}
		if params.Type != "" {
			query.Set("type", params.Type)
		}
		if len(query) > 0 {
			urlPath += "?" + query.Encode()
		}
	}

	var response ListDevicesResponse
	err := c.do(ctx, http.MethodGet, urlPath, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list devices: %w", err)
	}

	return &response, nil
}

// GetDevice retrieves a specific device by ID
func (c *Client) GetDevice(ctx context.Context, siteID, deviceID string) (*Device, error) {
	var response struct {
		Data []Device `json:"data"`
	}

	err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v1/sites/%s/devices/%s", siteID, deviceID), nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get device: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("device not found: %s", deviceID)
	}

	return &response.Data[0], nil
}

// ExecutePortAction performs an action on a specific port of a device
func (c *Client) ExecutePortAction(ctx context.Context, siteID, deviceID string, action *DevicePortAction) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	urlPath := fmt.Sprintf("/api/v1/sites/%s/devices/%s/port/%s", siteID, deviceID, action.PortID)
	err := c.do(ctx, http.MethodPost, urlPath, action, nil)
	if err != nil {
		return fmt.Errorf("failed to execute port action: %w", err)
	}

	return nil
}

// ExecuteDeviceAction performs an action on a device
func (c *Client) ExecuteDeviceAction(ctx context.Context, siteID, deviceID string, action *DeviceAction) error {
	if action == nil {
		return fmt.Errorf("action cannot be nil")
	}

	urlPath := fmt.Sprintf("/api/v1/sites/%s/devices/%s", siteID, deviceID)
	err := c.do(ctx, http.MethodPost, urlPath, action, nil)
	if err != nil {
		return fmt.Errorf("failed to execute device action: %w", err)
	}

	return nil
}

// GetDeviceStatistics retrieves the latest statistics for a device
func (c *Client) GetDeviceStatistics(ctx context.Context, siteID, deviceID string) (*DeviceStatistics, error) {
	var response struct {
		Data []DeviceStatistics `json:"data"`
	}

	urlPath := fmt.Sprintf("/api/v1/sites/%s/devices/%s/stats", siteID, deviceID)
	err := c.do(ctx, http.MethodGet, urlPath, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get device statistics: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("no statistics found for device: %s", deviceID)
	}

	return &response.Data[0], nil
}
