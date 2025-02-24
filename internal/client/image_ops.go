package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

// CreateImage creates a new workspace image
func (c *Client) CreateImage(image *Image) (*Image, error) {
	// Convert Image to CreateImageRequest
	req := &CreateImageRequest{
		ImageSrc:           image.ImageSrc,
		Name:               image.Name,
		FriendlyName:       image.FriendlyName,
		Description:        image.Description,
		Memory:             image.Memory,
		Cores:              image.Cores,
		DockerRegistry:     image.DockerRegistry,
		UncompressedSizeMB: image.UncompressedSizeMB,
		ImageType:          image.ImageType,
		Enabled:            image.Enabled,
	}
	return c.AddWorkspaceImage(req)
}

func (c *Client) GetImage(imageID string) (*Image, error) {
	var image *Image
	err := c.retryOperation(func() error {
		payload := map[string]interface{}{
			"api_key":        c.APIKey,
			"api_key_secret": c.APISecret,
			"image_id":       imageID,
		}

		resp, err := c.doRequestLegacy("POST", "/api/public/get_image", payload)
		if err != nil {
			return fmt.Errorf("error making request: %v", err)
		}
		defer resp.Body.Close()

		if err := c.handleAPIError(resp, "get image"); err != nil {
			return err
		}

		var result struct {
			Image *Image `json:"image"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return fmt.Errorf("error decoding response: %v", err)
		}

		image = result.Image
		return nil
	})

	if err != nil {
		return nil, err
	}

	return image, nil
}

// UpdateImage updates an existing workspace image
func (c *Client) UpdateImage(image *Image) (*Image, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_image":   image,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/update_image",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Image        *Image `json:"image"`
		ErrorMessage string `json:"error_message,omitempty"`
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if result.ErrorMessage != "" {
		return nil, fmt.Errorf("API returned error: %s", result.ErrorMessage)
	}

	return result.Image, nil
}

// DeleteImage deletes a workspace image
func (c *Client) DeleteImage(imageID string) error {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_image": map[string]string{
			"image_id": imageID,
		},
	})
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/delete_image",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// GetImages retrieves all available images
func (c *Client) GetImages() ([]Image, error) {
	log.Printf("[DEBUG] Getting images from Kasm API")

	req := struct {
		APIKey       string `json:"api_key"`
		APIKeySecret string `json:"api_key_secret"`
	}{
		APIKey:       c.APIKey,
		APIKeySecret: c.APISecret,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	log.Printf("[DEBUG] GetImages request body: %s", string(reqBody))

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/get_images",
		"application/json",
		bytes.NewReader(reqBody),
	)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	log.Printf("[DEBUG] GetImages response status: %d", resp.StatusCode)
	log.Printf("[DEBUG] GetImages response body: %s", bodyString)

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, bodyString)
	}

	var result struct {
		Images []Image `json:"images"`
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, bodyString)
	}

	log.Printf("[DEBUG] Found %d images", len(result.Images))
	for i, img := range result.Images {
		log.Printf("[DEBUG] Image %d: ID=%s, Name=%s", i, img.ImageID, img.Name)
	}

	return result.Images, nil
}

// GetSessionRecordings retrieves recordings for a specific session
func (c *Client) GetSessionRecordings(kasmID string, preauthDownloadLink bool) ([]SessionRecording, error) {
	req := struct {
		APIKey              string `json:"api_key"`
		APIKeySecret        string `json:"api_key_secret"`
		TargetKasmID        string `json:"target_kasm_id"`
		PreauthDownloadLink bool   `json:"preauth_download_link"`
	}{
		APIKey:              c.APIKey,
		APIKeySecret:        c.APISecret,
		TargetKasmID:        kasmID,
		PreauthDownloadLink: preauthDownloadLink,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/get_session_recordings",
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
		SessionRecordings []SessionRecording `json:"session_recordings"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.SessionRecordings, nil
}

// GetSessionsRecordings retrieves recordings for multiple sessions
func (c *Client) GetSessionsRecordings(kasmIDs []string, preauthDownloadLink bool) (map[string][]SessionRecording, error) {
	req := struct {
		APIKey              string   `json:"api_key"`
		APIKeySecret        string   `json:"api_key_secret"`
		TargetKasmIDs       []string `json:"target_kasm_ids"`
		PreauthDownloadLink bool     `json:"preauth_download_link"`
	}{
		APIKey:              c.APIKey,
		APIKeySecret:        c.APISecret,
		TargetKasmIDs:       kasmIDs,
		PreauthDownloadLink: preauthDownloadLink,
	}

	reqBody, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/get_sessions_recordings",
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
		KasmSessions map[string]struct {
			SessionRecordings []SessionRecording `json:"session_recordings"`
		} `json:"kasm_sessions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// Convert to simpler map structure
	recordings := make(map[string][]SessionRecording)
	for kasmID, session := range result.KasmSessions {
		recordings[kasmID] = session.SessionRecordings
	}

	return recordings, nil
}

// AddWorkspaceImage adds a new workspace image to Kasm
func (c *Client) AddWorkspaceImage(image *CreateImageRequest) (*Image, error) {
	reqBody := createImageAPIRequest{
		APIKey:       c.APIKey,
		APIKeySecret: c.APISecret,
		TargetImage:  *image,
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	// Add debug logging for request
	if os.Getenv("KASM_DEBUG") != "" {
		log.Printf("[DEBUG] AddWorkspaceImage request body: %s", string(body))
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/create_image",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read the full response body
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %v", err)
	}

	// Add debug logging for response
	if os.Getenv("KASM_DEBUG") != "" {
		log.Printf("[DEBUG] AddWorkspaceImage response: %s", string(bodyBytes))
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Image        *Image `json:"image"`
		ErrorMessage string `json:"error_message,omitempty"`
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, string(bodyBytes))
	}

	if result.ErrorMessage != "" {
		return nil, fmt.Errorf("API returned error: %s", result.ErrorMessage)
	}

	if result.Image == nil {
		return nil, fmt.Errorf("API returned success but image is nil. Response body: %s", string(bodyBytes))
	}

	if result.Image.ImageID == "" {
		return nil, fmt.Errorf("API returned image but ImageID is empty. Response body: %s", string(bodyBytes))
	}

	return result.Image, nil
}
