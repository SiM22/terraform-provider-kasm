package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// GetLoginURL generates a login URL for the specified user
func (c *Client) GetLoginURL(userID string) (*GetLoginResponse, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user": map[string]string{
			"user_id": userID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_login", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result GetLoginResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
