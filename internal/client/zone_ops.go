package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetZones retrieves all deployment zones
func (c *Client) GetZones(brief bool) ([]Zone, error) {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"brief":          brief,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/get_zones", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result struct {
		Zones []Zone `json:"zones"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, string(body))
	}

	return result.Zones, nil
}

// GetZone retrieves a specific deployment zone by ID
func (c *Client) GetZone(zoneID string) (*Zone, error) {
	zones, err := c.GetZones(false)
	if err != nil {
		return nil, fmt.Errorf("error getting zones: %v", err)
	}

	for _, zone := range zones {
		if zone.ZoneID == zoneID {
			return &zone, nil
		}
	}

	return nil, fmt.Errorf("zone with ID %s not found", zoneID)
}
