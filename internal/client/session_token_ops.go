package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateSessionToken creates a new session token for a user
func (c *Client) CreateSessionToken(request *CreateSessionTokenRequest) (*SessionToken, error) {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user":    request.TargetUser,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/create_session_token", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		SessionToken *SessionToken `json:"session_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.SessionToken, nil
}

// GetSessionToken retrieves a specific session token
func (c *Client) GetSessionToken(request *GetSessionTokenRequest) (*SessionToken, error) {
	payload := map[string]interface{}{
		"api_key":              c.APIKey,
		"api_key_secret":       c.APISecret,
		"target_session_token": request.TargetSessionToken,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/get_session_token", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		SessionToken *SessionToken `json:"session_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.SessionToken, nil
}

// GetSessionTokens retrieves all session tokens for a user
func (c *Client) GetSessionTokens(request *GetSessionTokensRequest) ([]SessionToken, error) {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user":    request.TargetUser,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/get_session_tokens", payload)
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
		SessionTokens []SessionToken `json:"session_tokens"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, string(body))
	}

	return result.SessionTokens, nil
}

// UpdateSessionToken updates (promotes) an existing session token
func (c *Client) UpdateSessionToken(request *UpdateSessionTokenRequest) (*SessionToken, error) {
	payload := map[string]interface{}{
		"api_key":              c.APIKey,
		"api_key_secret":       c.APISecret,
		"target_session_token": request.TargetSessionToken,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/update_session_token", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		SessionToken *SessionToken `json:"session_token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.SessionToken, nil
}

// DeleteSessionToken deletes a specific session token
func (c *Client) DeleteSessionToken(request *DeleteSessionTokenRequest) error {
	payload := map[string]interface{}{
		"api_key":              c.APIKey,
		"api_key_secret":       c.APISecret,
		"target_session_token": request.TargetSessionToken,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/delete_session_token", payload)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete session token, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// DeleteSessionTokens deletes all session tokens for a user
func (c *Client) DeleteSessionTokens(request *DeleteSessionTokensRequest) error {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user":    request.TargetUser,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/delete_session_tokens", payload)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete session tokens, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}
