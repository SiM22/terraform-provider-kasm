package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

const (
	defaultRetryTimeout = 2 * time.Minute
	defaultRetryDelay   = 5 * time.Second
	maxRetries          = 3
)

func isRetryableError(err error) bool {
	if err == nil {
		return false
	}
	// Add conditions for retryable errors
	var apiErr *APIError
	if errors.As(err, &apiErr) {
		switch apiErr.StatusCode {
		case 429, 502, 503, 504:
			return true
		}
	}
	return false
}

// CreateCastConfig creates a new casting configuration
func (c *Client) CreateCastConfig(config *CastConfig) (*CastConfig, error) {
	req := createCastConfigRequest{
		APIKey:           c.APIKey,
		APIKeySecret:     c.APISecret,
		TargetCastConfig: *config,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/create_cast_config",
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		// Read error response
		var errorResp struct {
			ErrorMessage string `json:"error_message"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&errorResp); err == nil && errorResp.ErrorMessage != "" {
			return nil, fmt.Errorf("API request failed: %s", errorResp.ErrorMessage)
		}
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var result createCastConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result.CastConfig, nil
}

// GetCastConfig retrieves a casting configuration by ID
func (c *Client) GetCastConfig(id string) (*CastConfig, error) {
	req := getCastConfigRequest{
		APIKey:       c.APIKey,
		APIKeySecret: c.APISecret,
		CastConfigID: id,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/get_cast_config",
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var result getCastConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result.CastConfig, nil
}

// UpdateCastConfig updates an existing casting configuration
func (c *Client) UpdateCastConfig(config *CastConfig) (*CastConfig, error) {
	req := updateCastConfigRequest{
		APIKey:           c.APIKey,
		APIKeySecret:     c.APISecret,
		TargetCastConfig: *config,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/update_cast_config",
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var result updateCastConfigResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result.CastConfig, nil
}

// DeleteCastConfig deletes a casting configuration
func (c *Client) DeleteCastConfig(id string) error {
	req := deleteCastConfigRequest{
		APIKey:       c.APIKey,
		APIKeySecret: c.APISecret,
		CastConfigID: id,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/delete_cast_config",
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// ListCastConfigs retrieves all casting configurations
func (c *Client) ListCastConfigs() ([]CastConfig, error) {
	req := struct {
		APIKey       string `json:"api_key"`
		APIKeySecret string `json:"api_key_secret"`
	}{
		APIKey:       c.APIKey,
		APIKeySecret: c.APISecret,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/get_cast_configs",
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		CastConfigs []CastConfig `json:"cast_configs"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.CastConfigs, nil
}
