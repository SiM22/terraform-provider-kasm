package client

import (
	"encoding/json"
	"fmt"
)

// SetSessionPermissions sets permissions for multiple users in a session
func (c *Client) SetSessionPermissions(request *SetSessionPermissionsRequest) ([]SessionPermission, error) {
	payload := map[string]interface{}{
		"api_key":                    c.APIKey,
		"api_key_secret":             c.APISecret,
		"target_session_permissions": request.TargetSessionPermissions,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/set_session_permissions", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	// After setting permissions, get the user details for each permission
	var permissions []SessionPermission
	for _, perm := range request.TargetSessionPermissions.SessionPermissions {
		// Get user details
		user, err := c.GetUser(perm.UserID)
		if err != nil {
			return nil, fmt.Errorf("error getting user details: %v", err)
		}

		permissions = append(permissions, SessionPermission{
			UserID:      perm.UserID,
			Access:      perm.Access,
			Username:    user.Username,
			VNCUsername: fmt.Sprintf("vnc-%s", user.Username), // VNC username is typically prefixed
		})
	}

	return permissions, nil
}

// GetSessionPermissions retrieves session permissions for a session
func (c *Client) GetSessionPermissions(request *GetSessionPermissionsRequest) ([]SessionPermission, error) {
	payload := map[string]interface{}{
		"api_key":                    c.APIKey,
		"api_key_secret":             c.APISecret,
		"target_session_permissions": request.TargetSessionPermissions,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/get_session_permissions", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		SessionPermissions []SessionPermission `json:"session_permissions"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// If we got permissions but they don't have usernames, fetch them
	if len(result.SessionPermissions) > 0 {
		var permissions []SessionPermission
		for _, perm := range result.SessionPermissions {
			// Get user details
			user, err := c.GetUser(perm.UserID)
			if err != nil {
				return nil, fmt.Errorf("error getting user details: %v", err)
			}

			// Create new permission with original access level and user details
			permissions = append(permissions, SessionPermission{
				UserID:      perm.UserID,
				Access:      perm.Access, // Preserve the original access level
				Username:    user.Username,
				VNCUsername: fmt.Sprintf("vnc-%s", user.Username),
			})
		}
		return permissions, nil
	}

	return result.SessionPermissions, nil
}

// DeleteAllSessionPermissions deletes all session permissions for a session
func (c *Client) DeleteAllSessionPermissions(request *DeleteAllSessionPermissionsRequest) error {
	payload := map[string]interface{}{
		"api_key":                    c.APIKey,
		"api_key_secret":             c.APISecret,
		"target_session_permissions": request.TargetSessionPermissions,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/delete_all_session_permissions", payload)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
