package unifi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// HotspotVoucher represents a UniFi hotspot voucher
type HotspotVoucher struct {
	ID                  string `json:"_id"`                            // Unique identifier
	CreatedAt           string `json:"createdAt"`                      // Timestamp when the voucher was created
	Name                string `json:"name"`                           // Voucher note, may contain duplicate values across vouchers
	Code                string `json:"code"`                           // Secret code to activate the voucher using the Hotspot portal
	AuthorizeGuestLimit int    `json:"authorizedGuestLimit,omitempty"` // Optional limit for how many different guests can use the voucher
	AuthorizeGuestCount int    `json:"authorizedGuestCount"`           // For how many guests the voucher has been used to authorize access
	ActivatedAt         string `json:"activatedAt,omitempty"`          // Optional timestamp when the voucher was activated (first guest)
	ExpiresAt           string `json:"expiresAt,omitempty"`            // Optional timestamp when the voucher will expire
	Expired             bool   `json:"expired"`                        // Whether the voucher has expired and can no longer be used
	TimeLimitMinutes    int    `json:"timeLimitMinutes"`               // How long the voucher will provide access since authorization
	DataUsageLimitMB    int    `json:"dataUsageLimitMBytes,omitempty"` // Optional data usage limit in megabytes
	RxRateLimitKbps     int    `json:"rxRateLimitKbps,omitempty"`      // Optional download rate limit in kilobits per second
	TxRateLimitKbps     int    `json:"txRateLimitKbps,omitempty"`      // Optional upload rate limit in kilobits per second
}

// ListHotspotVouchersParams contains parameters for listing hotspot vouchers
type ListHotspotVouchersParams struct {
	Offset int `json:"offset,omitempty"`
	Limit  int `json:"limit,omitempty"`
}

// ListHotspotVouchersResponse represents the response from listing hotspot vouchers
type ListHotspotVouchersResponse struct {
	PaginatedResponse
	Data []HotspotVoucher `json:"data"`
}

// CreateHotspotVoucherRequest represents the request to create a hotspot voucher
type CreateHotspotVoucherRequest struct {
	Name                string `json:"name"`                           // Required: Voucher note
	TimeLimitMinutes    int    `json:"timeLimitMinutes"`               // Required: How long the voucher will provide access
	AuthorizeGuestLimit int    `json:"authorizedGuestLimit,omitempty"` // Optional limit for number of guests
	DataUsageLimitMB    int    `json:"dataUsageLimitMBytes,omitempty"` // Optional data usage limit in MB
	RxRateLimitKbps     int    `json:"rxRateLimitKbps,omitempty"`      // Optional download rate limit
	TxRateLimitKbps     int    `json:"txRateLimitKbps,omitempty"`      // Optional upload rate limit
	Count               int    `json:"count,omitempty"`                // Number of vouchers to create
}

// CreateHotspotVoucherResponse represents the response from creating hotspot vouchers
type CreateHotspotVoucherResponse struct {
	Data []HotspotVoucher `json:"data"`
}

// GenerateHotspotVouchersRequest represents the request to generate hotspot vouchers
type GenerateHotspotVouchersRequest struct {
	Count               int    `json:"count"`                          // [1..10000] Number of vouchers to generate, default: 1
	Name                string `json:"name"`                           // Required: Voucher note, duplicated across all generated vouchers
	AuthorizeGuestLimit int    `json:"authorizedGuestLimit,omitempty"` // [1..] Optional limit for guests per voucher
	TimeLimitMinutes    int    `json:"timeLimitMinutes"`               // [1..1000000] Required: How long the voucher provides access
	DataUsageLimitMB    int    `json:"dataUsageLimitMBytes,omitempty"` // [1..1046576] Optional data usage limit in MB
	RxRateLimitKbps     int    `json:"rxRateLimitKbps,omitempty"`      // [2..100000] Optional download rate limit in Kbps
	TxRateLimitKbps     int    `json:"txRateLimitKbps,omitempty"`      // [2..100000] Optional upload rate limit in Kbps
}

// GenerateHotspotVouchersResponse represents the response from generating vouchers
type GenerateHotspotVouchersResponse struct {
	Data []HotspotVoucher `json:"data"`
}

// GetVoucherDetailsResponse represents the response from getting voucher details
type GetVoucherDetailsResponse struct {
	Data []HotspotVoucher `json:"data"`
}

// ListHotspotVouchers retrieves a paginated list of hotspot vouchers for a site
func (c *Client) ListHotspotVouchers(ctx context.Context, siteID string, params *ListHotspotVouchersParams) (*ListHotspotVouchersResponse, error) {
	urlPath := fmt.Sprintf("/api/v1/sites/%s/hotspot/vouchers", siteID)

	if params != nil {
		query := url.Values{}
		if params.Offset > 0 {
			query.Set("offset", fmt.Sprint(params.Offset))
		}
		if params.Limit > 0 {
			query.Set("limit", fmt.Sprint(params.Limit))
		}
		if len(query) > 0 {
			urlPath += "?" + query.Encode()
		}
	}

	var response ListHotspotVouchersResponse
	err := c.do(ctx, http.MethodGet, urlPath, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list hotspot vouchers: %w", err)
	}

	return &response, nil
}

// CreateHotspotVoucher creates one or more hotspot vouchers for a site
func (c *Client) CreateHotspotVoucher(ctx context.Context, siteID string, request *CreateHotspotVoucherRequest) (*CreateHotspotVoucherResponse, error) {
	urlPath := fmt.Sprintf("/api/v1/sites/%s/hotspot/vouchers", siteID)

	var response CreateHotspotVoucherResponse
	err := c.do(ctx, http.MethodPost, urlPath, request, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to create hotspot voucher: %w", err)
	}

	return &response, nil
}

// GetHotspotVoucher retrieves a specific hotspot voucher by ID
func (c *Client) GetHotspotVoucher(ctx context.Context, siteID, voucherID string) (*HotspotVoucher, error) {
	var response struct {
		Data []HotspotVoucher `json:"data"`
	}

	err := c.do(ctx, http.MethodGet, fmt.Sprintf("/api/v1/sites/%s/hotspot/vouchers/%s", siteID, voucherID), nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get hotspot voucher: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("voucher not found: %s", voucherID)
	}

	return &response.Data[0], nil
}

// DeleteHotspotVoucher deletes a specific hotspot voucher
func (c *Client) DeleteHotspotVoucher(ctx context.Context, siteID, voucherID string) error {
	err := c.do(ctx, http.MethodDelete, fmt.Sprintf("/api/v1/sites/%s/hotspot/vouchers/%s", siteID, voucherID), nil, nil)
	if err != nil {
		return fmt.Errorf("failed to delete hotspot voucher: %w", err)
	}

	return nil
}

// GenerateHotspotVouchers generates one or more hotspot vouchers with the specified parameters
func (c *Client) GenerateHotspotVouchers(ctx context.Context, siteID string, request *GenerateHotspotVouchersRequest) (*GenerateHotspotVouchersResponse, error) {
	if request == nil {
		return nil, fmt.Errorf("request cannot be nil")
	}

	// Validate required fields and ranges
	if request.Name == "" {
		return nil, fmt.Errorf("name is required")
	}
	if request.Count < 1 || request.Count > 10000 {
		return nil, fmt.Errorf("count must be between 1 and 10000")
	}
	if request.TimeLimitMinutes < 1 || request.TimeLimitMinutes > 1000000 {
		return nil, fmt.Errorf("timeLimitMinutes must be between 1 and 1000000")
	}
	if request.AuthorizeGuestLimit < 0 {
		return nil, fmt.Errorf("authorizedGuestLimit must be greater than 0")
	}
	if request.DataUsageLimitMB != 0 && (request.DataUsageLimitMB < 1 || request.DataUsageLimitMB > 1046576) {
		return nil, fmt.Errorf("dataUsageLimitMBytes must be between 1 and 1046576")
	}
	if request.RxRateLimitKbps != 0 && (request.RxRateLimitKbps < 2 || request.RxRateLimitKbps > 100000) {
		return nil, fmt.Errorf("rxRateLimitKbps must be between 2 and 100000")
	}
	if request.TxRateLimitKbps != 0 && (request.TxRateLimitKbps < 2 || request.TxRateLimitKbps > 100000) {
		return nil, fmt.Errorf("txRateLimitKbps must be between 2 and 100000")
	}

	urlPath := fmt.Sprintf("/api/v1/sites/%s/hotspot/vouchers/create", siteID)

	var response GenerateHotspotVouchersResponse
	err := c.do(ctx, http.MethodPost, urlPath, request, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to generate hotspot vouchers: %w", err)
	}

	return &response, nil
}

// GetVoucherDetails retrieves detailed information about a specific hotspot voucher
func (c *Client) GetVoucherDetails(ctx context.Context, siteID, voucherID string) (*HotspotVoucher, error) {
	if siteID == "" {
		return nil, fmt.Errorf("siteId is required")
	}
	if voucherID == "" {
		return nil, fmt.Errorf("voucherId is required")
	}

	urlPath := fmt.Sprintf("/api/v1/sites/%s/hotspot/vouchers/%s", siteID, voucherID)

	var response GetVoucherDetailsResponse
	err := c.do(ctx, http.MethodGet, urlPath, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get voucher details: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("voucher not found: %s", voucherID)
	}

	return &response.Data[0], nil
}
