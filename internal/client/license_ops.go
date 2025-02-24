package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// Activate activates a Kasm license using the provided activation key
func (c *Client) Activate(req *ActivateRequest) (*License, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"activation_key": req.ActivationKey,
		"seats":          req.Seats,
		"issued_to":      req.IssuedTo,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/activate", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result ActivateResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result.License, nil
}

// GetLicenses retrieves all Kasm licenses
func (c *Client) GetLicenses() ([]License, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_licenses", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result GetLicensesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if !result.Success {
		return nil, fmt.Errorf("API request was not successful: %s", result.Message)
	}

	return result.Data, nil
}
