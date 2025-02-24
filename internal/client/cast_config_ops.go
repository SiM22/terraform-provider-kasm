package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// GetCastingConfigs with retry logic
func (c *Client) GetCastingConfigs() ([]CastingConfig, error) {
	var configs []CastingConfig
	err := c.retryOperation(func() error {
		payload := map[string]interface{}{
			"api_key":        c.APIKey,
			"api_key_secret": c.APISecret,
		}

		resp, err := c.doRequestLegacy("POST", "/api/public/get_cast_configs", payload)
		if err != nil {
			return fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read error response: %v", err)
			}
			return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(body))
		}

		var result struct {
			CastConfigs []CastingConfig `json:"cast_configs"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("error decoding response: %v", err)
		}

		configs = result.CastConfigs
		return nil
	})

	if err != nil {
		return nil, err
	}

	return configs, nil
}

// GetCastingConfig retrieves a specific casting configuration
func (c *Client) GetCastingConfig(castConfigID string) (*CastingConfig, error) {
	var config *CastingConfig
	err := c.retryOperation(func() error {
		payload := map[string]interface{}{
			"api_key":        c.APIKey,
			"api_key_secret": c.APISecret,
			"cast_config_id": castConfigID,
		}

		resp, err := c.doRequestLegacy("POST", "/api/public/get_cast_config", payload)
		if err != nil {
			return fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read error response: %v", err)
			}
			return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(body))
		}

		var result struct {
			CastConfig *CastingConfig `json:"cast_config"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("error decoding response: %v", err)
		}

		config = result.CastConfig
		return nil
	})

	if err != nil {
		return nil, err
	}

	return config, nil
}

// CreateCastingConfig creates a new casting configuration
func (c *Client) CreateCastingConfig(request *CreateCastingConfigRequest) (*CastingConfig, error) {
	if err := c.validateCastingConfig(&request.TargetCastConfig); err != nil {
		return nil, err
	}

	var config *CastingConfig
	err := c.retryOperation(func() error {
		payload := map[string]interface{}{
			"api_key":            c.APIKey,
			"api_key_secret":     c.APISecret,
			"target_cast_config": request.TargetCastConfig,
		}

		resp, err := c.doRequestLegacy("POST", "/api/public/create_cast_config", payload)
		if err != nil {
			return fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read error response: %v", err)
			}
			return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(body))
		}

		var result struct {
			CastConfig *CastingConfig `json:"cast_config"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("error decoding response: %v", err)
		}

		config = result.CastConfig
		return nil
	})

	if err != nil {
		return nil, err
	}

	return config, nil
}

// UpdateCastingConfig updates an existing casting configuration
func (c *Client) UpdateCastingConfig(request *UpdateCastingConfigRequest) (*CastingConfig, error) {
	if err := c.validateCastingConfig(&request.TargetCastConfig); err != nil {
		return nil, err
	}

	var config *CastingConfig
	err := c.retryOperation(func() error {
		payload := map[string]interface{}{
			"api_key":            c.APIKey,
			"api_key_secret":     c.APISecret,
			"target_cast_config": request.TargetCastConfig,
		}

		resp, err := c.doRequestLegacy("POST", "/api/public/update_cast_config", payload)
		if err != nil {
			return fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		if err := c.handleAPIError(resp, "update cast config"); err != nil {
			return err
		}

		var result struct {
			CastConfig *CastingConfig `json:"cast_config"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("error decoding response: %v", err)
		}

		config = result.CastConfig
		return nil
	})

	if err != nil {
		return nil, err
	}

	return config, nil
}

// DeleteCastingConfig deletes a casting configuration
func (c *Client) DeleteCastingConfig(request *DeleteCastingConfigRequest) error {
	return c.retryOperation(func() error {
		payload := map[string]interface{}{
			"api_key":             c.APIKey,
			"api_key_secret":      c.APISecret,
			"cast_config_id":      request.CastConfigID,
			"casting_config_name": request.CastingConfigName,
		}

		resp, err := c.doRequestLegacy("POST", "/api/public/delete_cast_config", payload)
		if err != nil {
			return fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		if resp.StatusCode >= 400 {
			body, err := io.ReadAll(resp.Body)
			if err != nil {
				return fmt.Errorf("failed to read error response: %v", err)
			}
			return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(body))
		}
		return nil
	})
}

// ValidateCastingConfig validates the configuration before sending to the API
func (c *Client) ValidateCastingConfig(config *CastingConfig) error {
	if config == nil {
		return fmt.Errorf("cast config cannot be nil")
	}

	if config.CastingConfigName == "" {
		return fmt.Errorf("casting_config_name is required")
	}

	if config.ImageID == "" {
		return fmt.Errorf("image_id is required")
	}

	if config.Key == "" {
		return fmt.Errorf("key is required")
	}

	if config.LimitSessions && config.SessionRemaining <= 0 {
		return fmt.Errorf("session_remaining must be greater than 0 when limit_sessions is true")
	}

	if config.LimitIPs {
		if config.IPRequestLimit <= 0 {
			return fmt.Errorf("ip_request_limit must be greater than 0 when limit_ips is true")
		}
		if config.IPRequestSeconds <= 0 {
			return fmt.Errorf("ip_request_seconds must be greater than 0 when limit_ips is true")
		}
	}

	return nil
}

// GetCastingConfigByName retrieves a casting configuration by its name
func (c *Client) GetCastingConfigByName(name string) (*CastingConfig, error) {
	configs, err := c.GetCastingConfigs()
	if err != nil {
		return nil, fmt.Errorf("error getting cast configs: %v", err)
	}

	for _, config := range configs {
		if config.CastingConfigName == name {
			return &config, nil
		}
	}

	return nil, fmt.Errorf("cast config with name %s not found", name)
}

// GetCastingConfigByKey retrieves a casting configuration by its key
func (c *Client) GetCastingConfigByKey(key string) (*CastingConfig, error) {
	configs, err := c.GetCastingConfigs()
	if err != nil {
		return nil, fmt.Errorf("error getting cast configs: %v", err)
	}

	for _, config := range configs {
		if config.Key == key {
			return &config, nil
		}
	}

	return nil, fmt.Errorf("cast config with key %s not found", key)
}

// validateCastingConfig validates the configuration before API operations
func (c *Client) validateCastingConfig(config *CastingConfig) error {
	if config == nil {
		return fmt.Errorf("cast config cannot be nil")
	}

	if config.CastingConfigName == "" {
		return fmt.Errorf("casting_config_name is required")
	}

	if config.ImageID == "" {
		return fmt.Errorf("image_id is required")
	}

	if config.Key == "" {
		return fmt.Errorf("key is required")
	}

	if config.LimitSessions && config.SessionRemaining <= 0 {
		return fmt.Errorf("session_remaining must be greater than 0 when limit_sessions is true")
	}

	if config.LimitIPs {
		if config.IPRequestLimit <= 0 {
			return fmt.Errorf("ip_request_limit must be greater than 0 when limit_ips is true")
		}
		if config.IPRequestSeconds <= 0 {
			return fmt.Errorf("ip_request_seconds must be greater than 0 when limit_ips is true")
		}
	}

	return nil
}

// handleAPIError processes API responses and returns appropriate errors
func (c *Client) handleAPIError(resp *http.Response, operation string) error {
	if resp.StatusCode >= 400 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to read error response: %v", err)
		}
		return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

// retryOperation retries an operation with exponential backoff
func (c *Client) retryOperation(operation func() error) error {
	maxRetries := 3
	backoff := 1 * time.Second

	var lastErr error
	for i := 0; i < maxRetries; i++ {
		if err := operation(); err != nil {
			lastErr = err
			if i < maxRetries-1 {
				time.Sleep(backoff)
				backoff *= 2
				continue
			}
		} else {
			return nil
		}
	}
	return fmt.Errorf("operation failed after %d retries: %v", maxRetries, lastErr)
}

// GetCastingConfigsByImage retrieves all casting configurations for a specific image
func (c *Client) GetCastingConfigsByImage(imageID string) ([]CastingConfig, error) {
	configs, err := c.GetCastingConfigs()
	if err != nil {
		return nil, err
	}

	var filtered []CastingConfig
	for _, config := range configs {
		if config.ImageID == imageID {
			filtered = append(filtered, config)
		}
	}
	return filtered, nil
}

// GetCastingConfigsByGroup retrieves all casting configurations for a specific group
func (c *Client) GetCastingConfigsByGroup(groupID string) ([]CastingConfig, error) {
	configs, err := c.GetCastingConfigs()
	if err != nil {
		return nil, err
	}

	var filtered []CastingConfig
	for _, config := range configs {
		if config.GroupID != nil && *config.GroupID == groupID {
			filtered = append(filtered, config)
		}
	}
	return filtered, nil
}
