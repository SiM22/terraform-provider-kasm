package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

// FrameStatsRequest represents the request body for the frame_stats API endpoint
type FrameStatsRequest struct {
	APIKey    string `json:"api_key"`
	APISecret string `json:"api_key_secret"`
	KasmID    string `json:"kasm_id"`
	UserID    string `json:"user_id,omitempty"`
	Client    string `json:"client,omitempty"`
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
		Client:    "auto", // Default to auto as per documentation
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	log.Printf("[DEBUG] Getting frame stats for kasm %s", kasmID)

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_kasm_frame_stats", "application/json", bytes.NewBuffer(body))
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

	// Check for error message in the response
	if result.ErrorMessage != "" {
		log.Printf("[DEBUG] Frame stats API returned error: %s", result.ErrorMessage)
		if result.ErrorMessage == "Error retrieving frame stats with status code (503)" ||
			result.ErrorMessage == "Error retrieving frame stats with status code (502)" {
			return nil, fmt.Errorf("frame stats unavailable: a user must be actively connected to the session")
		}
		return nil, fmt.Errorf("frame stats API error: %s", result.ErrorMessage)
	}

	return &result, nil
}
