package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

func (c *Client) CreateUser(user *User) (*User, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user": map[string]interface{}{
			"username":          user.Username,
			"first_name":        user.FirstName,
			"last_name":         user.LastName,
			"password":          user.Password,
			"locked":            user.Locked,
			"disabled":          user.Disabled,
			"organization":      user.Organization,
			"phone":             user.Phone,
			"authorized_images": user.AuthorizedImages,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(
		c.BaseURL+"/api/public/create_user",
		"application/json",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, response: %s",
			resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		User         User   `json:"user"`
		ErrorMessage string `json:"error_message"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	if result.ErrorMessage != "" {
		return nil, fmt.Errorf("API returned error: %s", result.ErrorMessage)
	}

	return &result.User, nil
}

func (c *Client) GetUser(userID string) (*User, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user": map[string]string{
			"user_id": userID,
		},
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_user", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		User User `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.User, nil
}

func (c *Client) UpdateUser(user *User) (*User, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user": map[string]interface{}{
			"user_id":           user.UserID,
			"username":          user.Username,
			"first_name":        user.FirstName,
			"last_name":         user.LastName,
			"password":          user.Password,
			"locked":            user.Locked,
			"disabled":          user.Disabled,
			"organization":      user.Organization,
			"phone":             user.Phone,
			"authorized_images": user.AuthorizedImages,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/update_user", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		User User `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result.User, nil
}

func (c *Client) DeleteUser(userID string) error {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user": map[string]interface{}{
			"user_id": userID,
		},
		"force": true,
	})
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/delete_user", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	return nil
}

func (c *Client) LogoutUser(userID string) error {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user": map[string]string{
			"user_id": userID,
		},
	})
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/logout_user", "application/json", bytes.NewBuffer(body))
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

func (c *Client) GetUserAttributes(userID string) (*UserAttributes, error) {
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

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_attributes", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		UserAttributes UserAttributes `json:"user_attributes"`
	}
	if err := json.NewDecoder(bytes.NewReader(bodyBytes)).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return &result.UserAttributes, nil
}

func (c *Client) UpdateUserAttributes(userID string, attributes map[string]interface{}) error {
	userAttrs := UserAttributes{
		UserID: userID,
	}

	// Map the attributes to their proper fields
	for k, v := range attributes {
		if str, ok := v.(string); ok {
			switch k {
			case "theme":
				userAttrs.Theme = str
			case "language":
				userAttrs.PreferredLanguage = str
			}
		}
	}

	body, err := json.Marshal(map[string]interface{}{
		"api_key":                c.APIKey,
		"api_key_secret":         c.APISecret,
		"target_user_attributes": userAttrs,
	})
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/update_user_attributes", "application/json", bytes.NewBuffer(body))
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

func (c *Client) GetUsers() ([]User, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"page":           0,
		"page_size":      100, // Get a reasonable number of users
		"sort_by":        "username",
		"sort_direction": "asc",
		"anonymous":      false,
		"anonymous_only": false,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_users", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Users []User `json:"users"`
		Total int    `json:"total"`
		Page  int    `json:"page"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Users, nil
}

// waitForUserGroups waits for the user's group memberships to be fully processed
func (c *Client) waitForUserGroups(userID string, expectedGroups []string) error {
	maxAttempts := 10
	delay := 2 * time.Second

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		user, err := c.GetUser(userID)
		if err != nil {
			return fmt.Errorf("error getting user: %v", err)
		}

		// Create a map of expected group names for easier lookup
		expectedGroupMap := make(map[string]bool)
		for _, groupName := range expectedGroups {
			expectedGroupMap[groupName] = true
		}

		// Check if all expected groups are present
		foundGroups := 0
		for _, group := range user.Groups {
			if expectedGroupMap[group.Name] {
				foundGroups++
			}
		}

		if foundGroups == len(expectedGroups) {
			return nil
		}

		if attempt < maxAttempts {
			time.Sleep(delay)
		}
	}

	return fmt.Errorf("timeout waiting for user groups to be processed")
}

func (c *Client) UpdateUserGroupsByName(userID string, groupNames []string) error {
	groups, err := c.GetGroups()
	if err != nil {
		return fmt.Errorf("error getting groups: %v", err)
	}

	// Create a map of group names to IDs
	groupMap := make(map[string]string)
	for _, group := range groups {
		groupMap[group.Name] = group.GroupID
	}

	// Get current user's groups
	currentUser, err := c.GetUser(userID)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	// Create map of current group IDs
	currentGroupIDs := make(map[string]bool)
	for _, group := range currentUser.Groups {
		currentGroupIDs[group.GroupID] = true
	}

	// Create map of desired group IDs
	desiredGroupIDs := make(map[string]bool)
	for _, name := range groupNames {
		id, ok := groupMap[name]
		if !ok {
			return fmt.Errorf("group not found: %s", name)
		}
		desiredGroupIDs[id] = true
	}

	// Remove user from groups they shouldn't be in
	for groupID := range currentGroupIDs {
		if !desiredGroupIDs[groupID] {
			err := c.RemoveUserFromGroup(userID, groupID)
			if err != nil {
				return fmt.Errorf("error removing user from group: %v", err)
			}
		}
	}

	// Add user to new groups
	for groupID := range desiredGroupIDs {
		if !currentGroupIDs[groupID] {
			err := c.AddUserToGroup(userID, groupID)
			if err != nil {
				return fmt.Errorf("error adding user to group: %v", err)
			}
		}
	}

	return nil
}

func (c *Client) GetUserAuthorizedImages(userID string) ([]string, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user": map[string]interface{}{
			"user_id": userID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_user", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		User User `json:"user"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.User.AuthorizedImages, nil
}

func (c *Client) UpdateUserAuthorizedImages(userID string, imageIDs []string) error {
	// First get the current user to preserve all fields
	currentUser, err := c.GetUser(userID)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}

	// Update only the authorized images field
	currentUser.AuthorizedImages = imageIDs

	// Use the update_user endpoint with all fields
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_user": map[string]interface{}{
			"user_id":           currentUser.UserID,
			"username":          currentUser.Username,
			"first_name":        currentUser.FirstName,
			"last_name":         currentUser.LastName,
			"organization":      currentUser.Organization,
			"phone":             currentUser.Phone,
			"locked":            currentUser.Locked,
			"disabled":          currentUser.Disabled,
			"authorized_images": imageIDs,
		},
	})
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/update_user", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// IsUserNotFoundError checks if the error is due to a user not being found
func IsUserNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "User not found")
}
