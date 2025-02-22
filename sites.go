package unifi

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
)

// Site represents a UniFi site
type Site struct {
	ID   string `json:"id"`   // Unique identifier
	Name string `json:"name"` // Site name
}

// ListSitesParams contains parameters for listing sites
type ListSitesParams struct {
	Offset int `json:"offset,omitempty"` // Default: 0
	Limit  int `json:"limit,omitempty"`  // [0..200] Default: 25
}

// ListSitesResponse represents the response from listing sites
type ListSitesResponse struct {
	Offset     int    `json:"offset"`     // Starting offset
	Limit      int    `json:"limit"`      // Number of sites per page
	Count      int    `json:"count"`      // Number of sites in this response
	TotalCount int    `json:"totalCount"` // Total number of sites available
	Data       []Site `json:"data"`       // List of sites
}

// ListSites retrieves all sites accessible to the authenticated user
// If Multi-Site option is enabled, returns all created sites.
// If Multi-Site option is disabled, returns just the default site.
func (c *Client) ListSites(ctx context.Context, params *ListSitesParams) (*ListSitesResponse, error) {
	const maxLimit = 200
	urlPath := "/v1/sites"

	if params != nil {
		query := url.Values{}
		if params.Offset > 0 {
			query.Set("offset", fmt.Sprint(params.Offset))
		}
		if params.Limit > 0 {
			if params.Limit > maxLimit {
				return nil, fmt.Errorf("limit must be between 0 and %d", maxLimit)
			}
			query.Set("limit", fmt.Sprint(params.Limit))
		}
		if len(query) > 0 {
			urlPath += "?" + query.Encode()
		}
	}

	var response ListSitesResponse
	if err := c.do(ctx, http.MethodGet, urlPath, nil, &response); err != nil {
		return nil, fmt.Errorf("failed to list sites: %w", err)
	}

	return &response, nil
}
