package unifi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Site represents a UniFi site
type Site struct {
	ID          string `json:"_id"`            // Unique identifier
	Name        string `json:"name"`           // Site name
	Description string `json:"desc"`           // Site description
	Role        string `json:"role"`           // User's role at this site
	Hidden      bool   `json:"attr_hidden"`    // Whether the site is hidden in the UI
	NoDelete    bool   `json:"attr_no_delete"` // Whether the site can be deleted
}

// ListSitesParams contains parameters for listing sites
type ListSitesParams struct {
	Offset int `json:"offset,omitempty"` // Default: 0
	Limit  int `json:"limit,omitempty"`  // [0..200] Default: 25
}

// ListSitesResponse represents the response from listing sites
type ListSitesResponse struct {
	PaginatedResponse
	Data []Site `json:"data"`
}

// ListSites retrieves all sites accessible to the authenticated user
// If Multi-Site option is enabled, returns all created sites.
// If Multi-Site option is disabled, returns just the default site.
func (c *Client) ListSites(ctx context.Context, params *ListSitesParams) (*ListSitesResponse, error) {
	urlPath := "/v1/sites"

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

	var response ListSitesResponse
	err := c.do(ctx, http.MethodGet, urlPath, nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to list sites: %w", err)
	}

	return &response, nil
}

// GetSite retrieves a specific site by ID
func (c *Client) GetSite(ctx context.Context, siteID string) (*Site, error) {
	if siteID == "" {
		return nil, fmt.Errorf("siteID cannot be empty")
	}

	var response struct {
		Data []Site `json:"data"`
	}

	err := c.do(ctx, http.MethodGet, fmt.Sprintf("/v1/sites/%s", siteID), nil, &response)
	if err != nil {
		return nil, fmt.Errorf("failed to get site: %w", err)
	}

	if len(response.Data) == 0 {
		return nil, fmt.Errorf("site not found: %s", siteID)
	}

	return &response.Data[0], nil
}
