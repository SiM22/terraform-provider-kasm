package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetStagingConfigs retrieves all staging configurations
func (c *Client) GetStagingConfigs() ([]StagingConfig, error) {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/get_staging_configs", payload)
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
		StagingConfigs []StagingConfig `json:"staging_configs"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, string(body))
	}
	return result.StagingConfigs, nil
}

// GetStagingConfig retrieves a specific staging configuration
func (c *Client) GetStagingConfig(stagingConfigID string) (*StagingConfig, error) {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_staging_config": map[string]string{
			"staging_config_id": stagingConfigID,
		},
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/get_staging_config", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		StagingConfig *StagingConfig `json:"staging_config"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.StagingConfig, nil
}

// CreateStagingConfig creates a new staging configuration
func (c *Client) CreateStagingConfig(request *CreateStagingConfigRequest) (*StagingConfig, error) {
	payload := map[string]interface{}{
		"api_key":               c.APIKey,
		"api_key_secret":        c.APISecret,
		"target_staging_config": request,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/create_staging_config", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		StagingConfig *StagingConfig `json:"staging_config"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.StagingConfig, nil
}

// UpdateStagingConfig updates an existing staging configuration
func (c *Client) UpdateStagingConfig(request *UpdateStagingConfigRequest) (*StagingConfig, error) {
	payload := map[string]interface{}{
		"api_key":               c.APIKey,
		"api_key_secret":        c.APISecret,
		"target_staging_config": request,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/update_staging_config", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		StagingConfig *StagingConfig `json:"staging_config"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.StagingConfig, nil
}

// DeleteStagingConfig deletes a staging configuration
func (c *Client) DeleteStagingConfig(stagingConfigID string) error {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_staging_config": map[string]string{
			"staging_config_id": stagingConfigID,
		},
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/delete_staging_config", payload)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete staging config, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
