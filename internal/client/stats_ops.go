package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// FrameStatsRequest represents the request body for the frame_stats API endpoint
type FrameStatsRequest struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_key_secret"`
	KasmID    string `json:"kasm_id"`
	UserID    string `json:"user_id,omitempty"`
}

// GetFrameStats retrieves frame statistics for a specific Kasm session
// It takes a kasmID parameter which is the ID of the Kasm session to get stats for
// and an optional userID parameter which is the ID of the user who owns the session.
// Returns a FrameStatsResponse containing the frame statistics, or an error if the request fails.
func (c *Client) GetFrameStats(kasmID string, userID string) (*FrameStatsResponse, error) {
	requestBody := FrameStatsRequest{
		APIKey:    c.APIKey,
		APISecret: c.APISecret,
		KasmID:    kasmID,
		UserID:    userID,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/frame_stats", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result FrameStatsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}
