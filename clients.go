package unifi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// NetworkClient represents a connected client device per the UniFi API
type NetworkClient struct {
	ID             string  `json:"id"`             // Unique identifier
	Name           string  `json:"name"`           // Client name
	ConnectedAt    string  `json:"connectedAt"`    // Connection timestamp
	IPAddress      string  `json:"ipAddress"`      // IP address
	Type           string  `json:"type"`           // Connection type (WIRED, WIRELESS, VPN)
	MACAddress     string  `json:"macAddress"`     // MAC address
	UplinkDeviceID string  `json:"uplinkDeviceId"` // ID of the device this client is connected to
	SiteID         string  `json:"site_id"`        // Site identifier
	Network        string  `json:"network"`        // Network name
	NetworkName    string  `json:"network_name"`   // Network display name
	OUI            string  `json:"oui"`            // Organizationally Unique Identifier
	LastSeen       int64   `json:"last_seen"`      // Last seen timestamp
	Uptime         int64   `json:"uptime"`         // Connection uptime in seconds
	IsWired        bool    `json:"is_wired"`       // Whether client is connected via wire
	IsGuest        bool    `json:"is_guest"`       // Whether client is on guest network
	DeviceID       string  `json:"device_id"`      // Connected device ID
	DeviceName     string  `json:"device_name"`    // Connected device name
	DeviceMAC      string  `json:"device_mac"`     // Connected device MAC
	RxBytes        int64   `json:"rx_bytes"`       // Received bytes
	TxBytes        int64   `json:"tx_bytes"`       // Transmitted bytek
	RxRate         float64 `json:"rx_rate"`        // Current receive rate
	TxRate         float64 `json:"tx_rate"`        // Current transmit rate
	SignalStrength int     `json:"signal"`         // Signal strength (for wireless)
	NoiseFloor     int     `json:"noise"`          // Noise floor (for wireless)
	SNR            int     `json:"snr"`            // Signal to noise ratio (for wireless)
	Channel        int     `json:"channel"`        // Wireless channel
	RadioProtocol  string  `json:"radio_proto"`    // Radio protocol
	RadioBand      string  `json:"radio"`          // Radio band
	SSID           string  `json:"essid"`          // Connected SSID (for wireless)
	BSSID          string  `json:"bssid"`          // Connected BSSID (for wireless)
	UseFixedIP     bool    `json:"use_fixedip"`    // Whether using fixed IP
	FixedIP        string  `json:"fixed_ip"`       // Fixed IP address if set
	NetworkID      string  `json:"network_id"`     // Network identifier
}

// ListNetworkClientsParams contains parameters for listing network clients
type ListNetworkClientsParams struct {
	Offset int `json:"offset,omitempty"` // Default: 0
	Limit  int `json:"limit,omitempty"`  // [0..200] Default: 25
}

// ListNetworkClientsResponse represents the response from listing network clients
type ListNetworkClientsResponse struct {
	Offset     int             `json:"offset"`
	Limit      int             `json:"limit"`
	Count      int             `json:"count"`
	TotalCount int             `json:"totalCount"`
	Data       []NetworkClient `json:"data"`
}

// ListNetworkClients retrieves a paginated list of network clients for a site
func (c *Client) ListNetworkClients(ctx context.Context, siteID string, params *ListNetworkClientsParams) (*ListNetworkClientsResponse, error) {
	if siteID == "" {
		return nil, fmt.Errorf("siteId is required")
	}

	urlPath := fmt.Sprintf("/api/v1/sites/%s/clients", siteID)

	if params != nil {
		query := url.Values{}
		if params.Offset > 0 {
			query.Set("offset", fmt.Sprint(params.Offset))
		}
		if params.Limit > 0 {
			if params.Limit > 200 {
				return nil, fmt.Errorf("limit must be between 0 and 200")
			}
			query.Set("limit", fmt.Sprint(params.Limit))
		}
		if len(query) > 0 {
			urlPath += "?" + query.Encode()
		}
	}

	var response ListNetworkClientsResponse
	err := c.do(ctx, http.MethodGet, urlPath, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list network clients: %w", err)
	}

	return &response, nil
}

// GetNetworkClient retrieves a specific network client by ID
func (c *Client) GetNetworkClient(ctx context.Context, siteID, clientID string) (*NetworkClient, error) {
	if siteID == "" {
		return nil, fmt.Errorf("siteId is required")
	}
	if clientID == "" {
		return nil, fmt.Errorf("clientId is required")
	}

	var response struct {
		Data []NetworkClient `json:"data"`
	}

	err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v1/sites/%s/clients/%s", siteID, clientID), nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get network client: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("network client not found: %s", clientID)
	}

	return &response.Data[0], nil
}
