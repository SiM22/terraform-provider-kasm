package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func (c *Client) GetKasmStatus(userID, kasmID string, skipAgentCheck bool) (*KasmStatusResponse, error) {
	// Create request body
	requestBody := map[string]interface{}{
		"api_key":          c.APIKey,
		"api_key_secret":   c.APISecret,
		"user_id":          userID,
		"kasm_id":          kasmID,
		"skip_agent_check": skipAgentCheck,
	}

	// Marshal the request body
	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	// Log request details
	log.Printf("[DEBUG] GetKasmStatus request URL: %s", c.BaseURL+"/api/public/get_kasm_status")
	log.Printf("[DEBUG] GetKasmStatus request body: %s", string(body))
	log.Printf("[DEBUG] GetKasmStatus request for kasm_id: %s, user_id: %s", kasmID, userID)

	// Create a new request
	req, err := http.NewRequest("POST", c.BaseURL+"/api/public/get_kasm_status", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	log.Printf("[DEBUG] GetKasmStatus request headers: %v", req.Header)

	// Make the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)

	// Log response details
	log.Printf("[DEBUG] GetKasmStatus response status: %d", resp.StatusCode)
	log.Printf("[DEBUG] GetKasmStatus response headers: %v", resp.Header)
	log.Printf("[DEBUG] GetKasmStatus response body: %s", bodyString)

	// Handle different status codes
	if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
		return nil, fmt.Errorf("authentication error: %d, body: %s", resp.StatusCode, bodyString)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, bodyString)
	}

	// Decode the response
	var result KasmStatusResponse
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, bodyString)
	}

	// Check for API error messages
	if result.ErrorMessage != "" {
		log.Printf("[WARN] API returned error message: %s", result.ErrorMessage)
	}

	// Log successful response details
	if result.Kasm != nil {
		log.Printf("[DEBUG] Successfully retrieved status for kasm_id: %s, container_id: %s",
			kasmID, result.Kasm.ContainerID)
	} else {
		log.Printf("[DEBUG] Retrieved status for kasm_id: %s, but Kasm details are nil", kasmID)
	}

	return &result, nil
}

func (c *Client) JoinKasm(shareID string, userID string) (*JoinKasmResponse, error) {
	requestBody := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"share_id":       shareID,
		"user_id":        userID,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	log.Printf("[DEBUG] JoinKasm request URL: %s", c.BaseURL+"/api/public/join_kasm")
	log.Printf("[DEBUG] JoinKasm request body: %s", string(body))

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/join_kasm", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// Read response body for debugging
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	log.Printf("[DEBUG] JoinKasm response status: %d", resp.StatusCode)
	log.Printf("[DEBUG] JoinKasm response headers: %v", resp.Header)
	log.Printf("[DEBUG] JoinKasm response body: %s", bodyString)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, bodyString)
	}

	var result JoinKasmResponse
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, bodyString)
	}

	// If KasmURL is not set in the response, construct it from the base URL and kasm details
	if result.KasmURL == "" {
		result.KasmURL = fmt.Sprintf("/#/connect/kasm/%s/%s/%s",
			result.Kasm.KasmID,
			userID,
			result.SessionToken)
	}

	return &result, nil
}

func (c *Client) GetRDPConnectionInfo(userID, kasmID string, connectionType RDPConnectionType) (*RDPConnectionResponse, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":         c.APIKey,
		"api_key_secret":  c.APISecret,
		"user_id":         userID,
		"kasm_id":         kasmID,
		"connection_type": connectionType,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_rdp_client_connection_info", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result RDPConnectionResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

func (c *Client) GetKasms() (*GetKasmsResponse, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_kasms", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result GetKasmsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result, nil
}

func (c *Client) DestroyKasm(userID, kasmID string) error {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"user_id":        userID,
		"kasm_id":        kasmID,
	})
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/destroy_kasm", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// CreateKasmShareID creates a share ID for an existing Kasm session
// This uses an undocumented API endpoint that is used by the Kasm UI
func (c *Client) CreateKasmShareID(kasmID string, sessionToken string, username string) (string, error) {
	log.Printf("[DEBUG] Creating share ID for Kasm session %s", kasmID)

	// Create request body
	requestBody := map[string]interface{}{
		"kasm_id":  kasmID,
		"token":    sessionToken,
		"username": username,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request body: %v", err)
	}

	log.Printf("[DEBUG] CreateKasmShareID request URL: %s", c.BaseURL+"/api/create_kasm_share_id")
	log.Printf("[DEBUG] CreateKasmShareID request body: %s", string(body))

	// Create a new request
	req, err := http.NewRequest("POST", c.BaseURL+"/api/create_kasm_share_id", bytes.NewBuffer(body))
	if err != nil {
		return "", fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("error making request: %v", err)
	}

	// Read response body
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	resp.Body.Close()

	log.Printf("[DEBUG] CreateKasmShareID response status: %d", resp.StatusCode)
	log.Printf("[DEBUG] CreateKasmShareID response body: %s", bodyString)

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, bodyString)
	}

	// Parse the response
	var result struct {
		ShareID string `json:"share_id"`
	}
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
		return "", fmt.Errorf("error decoding response: %v, body: %s", err, bodyString)
	}

	if result.ShareID == "" {
		return "", fmt.Errorf("API returned empty share ID")
	}

	return result.ShareID, nil
}

// SetShareSettings enables or disables sharing for an existing Kasm session
func (c *Client) SetShareSettings(kasmID string, share bool, enableSharing bool) error {
	log.Printf("[DEBUG] Setting share settings for Kasm session %s: share=%v, enable_sharing=%v", kasmID, share, enableSharing)

	// Create request body
	requestBody := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"kasm_id":        kasmID,
		"share":          share,
		"enable_sharing": enableSharing,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	log.Printf("[DEBUG] SetShareSettings request URL: %s", c.BaseURL+"/api/public/set_share_settings")
	log.Printf("[DEBUG] SetShareSettings request body: %s", string(body))

	// Create a new request
	req, err := http.NewRequest("POST", c.BaseURL+"/api/public/set_share_settings", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error creating request: %v", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}

	// Read response body
	bodyBytes, _ := io.ReadAll(resp.Body)
	bodyString := string(bodyBytes)
	resp.Body.Close()

	log.Printf("[DEBUG] SetShareSettings response status: %d", resp.StatusCode)
	log.Printf("[DEBUG] SetShareSettings response body: %s", bodyString)

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, bodyString)
	}

	return nil
}

func (c *Client) CreateKasm(userID string, imageID string, sessionToken string, username string, share bool, persistent bool, allowResume bool, sessionAuthentication bool) (*CreateKasmResponse, error) {
	log.Printf("[DEBUG] Creating Kasm session for user %s with image %s", userID, imageID)

	// First check if the user has the image authorized
	user, err := c.GetUser(userID)
	if err != nil {
		return nil, fmt.Errorf("error getting user: %v", err)
	}

	// Check if the image is authorized through any of the user's groups
	isAuthorized := false
	for _, group := range user.Groups {
		groupImages, err := c.GetGroupImages(group.GroupID)
		if err != nil {
			log.Printf("[DEBUG] Error getting images for group %s: %v", group.GroupID, err)
			continue
		}
		for _, groupImage := range groupImages {
			if groupImage.ImageID == imageID {
				isAuthorized = true
				break
			}
		}
		if isAuthorized {
			break
		}
	}

	if !isAuthorized {
		log.Printf("[DEBUG] User %s is not authorized for image %s through any group", userID, imageID)
		return nil, fmt.Errorf("Image Not Authorized")
	}

	// First create a session token if not provided
	var token string
	if sessionToken == "" {
		createReq := &CreateSessionTokenRequest{}
		createReq.TargetUser.UserID = userID
		sessionTokenResp, err := c.CreateSessionToken(createReq)
		if err != nil {
			return nil, fmt.Errorf("error creating session token: %v", err)
		}
		token = sessionTokenResp.SessionToken
	} else {
		token = sessionToken
	}

	// Create request body according to API documentation
	// Always create session without sharing initially to avoid "Failed to create Share ID" error
	// Sharing will be enabled later through a separate API call if requested
	requestBody := map[string]interface{}{
		"api_key":                c.APIKey,
		"api_key_secret":         c.APISecret,
		"user_id":                userID,
		"image_id":               imageID,
		"share":                  false,
		"enable_sharing":         false,
		"environment":            map[string]string{},
		"session_token":          token,
		"persistent":             persistent,
		"allow_resume":           allowResume,
		"session_authentication": sessionAuthentication,
		"client_settings": map[string]interface{}{
			"allow_kasm_sharing": true, // Allow sharing in client settings even if not enabled initially
		},
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	log.Printf("[DEBUG] CreateKasm request URL: %s", c.BaseURL+"/api/public/request_kasm")
	log.Printf("[DEBUG] CreateKasm request body: %s", string(body))

	// Retry the request up to 3 times with exponential backoff
	var lastErr error
	for i := 0; i < 3; i++ {
		// Calculate backoff delay: 2^i * 1 second (1s, 2s, 4s)
		backoffDelay := time.Duration(1<<uint(i)) * time.Second

		// Create a new request
		req, err := http.NewRequest("POST", c.BaseURL+"/api/public/request_kasm", bytes.NewBuffer(body))
		if err != nil {
			return nil, fmt.Errorf("error creating request: %v", err)
		}

		// Set headers
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Accept", "application/json")

		log.Printf("[DEBUG] CreateKasm request headers: %v", req.Header)

		// Send the request
		resp, err := c.HTTPClient.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("error making request: %v", err)
			time.Sleep(backoffDelay)
			continue
		}

		// Read response body for debugging
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)
		resp.Body.Close()

		log.Printf("[DEBUG] CreateKasm response status: %d", resp.StatusCode)
		log.Printf("[DEBUG] CreateKasm response headers: %v", resp.Header)
		log.Printf("[DEBUG] CreateKasm response body: %s", bodyString)

		if resp.StatusCode == http.StatusUnauthorized || resp.StatusCode == http.StatusForbidden {
			lastErr = fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, bodyString)
			time.Sleep(backoffDelay)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, bodyString)
		}

		var result CreateKasmResponse
		if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
			return nil, fmt.Errorf("error decoding response: %v, body: %s", err, bodyString)
		}

		if result.ErrorMessage != "" {
			return nil, fmt.Errorf("API returned error: %s", result.ErrorMessage)
		}

		// If sharing was requested, enable it now that the session is created
		if share {
			log.Printf("[DEBUG] Enabling sharing for Kasm session %s", result.KasmID)
			err := c.SetShareSettings(result.KasmID, true, true)
			if err != nil {
				log.Printf("[WARN] Failed to enable sharing for Kasm session %s: %v", result.KasmID, err)
				// Don't fail the whole operation if enabling sharing fails
			} else {
				// Try to create a share ID using the undocumented API
				shareID, err := c.CreateKasmShareID(result.KasmID, token, username)
				if err != nil {
					log.Printf("[WARN] Failed to create share ID for Kasm session %s: %v", result.KasmID, err)
					// Don't fail the whole operation if creating share ID fails
				} else {
					// Update the result with the share ID
					result.ShareID = shareID
					log.Printf("[DEBUG] Successfully created share ID %s for Kasm session %s", shareID, result.KasmID)
				}
			}
		}

		return &result, nil
	}

	return nil, fmt.Errorf("failed after 3 retries: %v", lastErr)
}

// Keepalive sends a keepalive request to reset the expiration time of a Kasm session.
func (c *Client) Keepalive(kasmID string) (*KeepaliveResponse, error) {
	requestBody := KeepaliveRequest{
		APIKey:    c.APIKey,
		APISecret: c.APISecret,
		KasmID:    kasmID,
	}

	body, err := json.Marshal(requestBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	log.Printf("[DEBUG] Keepalive request URL: %s", c.BaseURL+"/api/public/keepalive")
	log.Printf("[DEBUG] Keepalive request body: %s", string(body))

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/keepalive", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("[WARN] Unexpected status code in keepalive: %d", resp.StatusCode)
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var keepaliveResponse KeepaliveResponse
	if err := json.NewDecoder(resp.Body).Decode(&keepaliveResponse); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &keepaliveResponse, nil
}
