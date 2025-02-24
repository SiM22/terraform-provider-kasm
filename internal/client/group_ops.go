package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// CreateGroup creates a new group
func (c *Client) CreateGroup(group *Group) (*Group, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_group": map[string]interface{}{
			"name":        group.Name,
			"priority":    group.Priority,
			"description": group.Description,
			"permissions": group.Permissions,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/create_group", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Group Group `json:"group"`
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	// If this group has sharing permissions, add the allow_kasm_sharing setting
	if group.Permissions != nil {
		for _, perm := range group.Permissions {
			if perm == "allow_kasm_sharing" {
				settingBody, err := json.Marshal(map[string]interface{}{
					"api_key":        c.APIKey,
					"api_key_secret": c.APISecret,
					"target_group": map[string]interface{}{
						"group_id": result.Group.GroupID,
					},
					"target_setting": map[string]interface{}{
						"group_setting_id": "13ef8423a5bb445cacbb6bd9f44a2454",
						"value":            "True",
					},
				})
				if err != nil {
					return nil, fmt.Errorf("error marshaling setting request: %v", err)
				}

				settingResp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/add_settings_group", "application/json", bytes.NewBuffer(settingBody))
				if err != nil {
					return nil, fmt.Errorf("error setting allow_kasm_sharing: %v", err)
				}
				settingResp.Body.Close()
				break
			}
		}
	}

	return &result.Group, nil
}

// GetGroup retrieves a group by ID
func (c *Client) GetGroup(groupID string) (*Group, error) {
	groups, err := c.GetGroups()
	if err != nil {
		return nil, fmt.Errorf("error getting groups: %v", err)
	}

	for _, group := range groups {
		if group.GroupID == groupID {
			return &group, nil
		}
	}

	return nil, &NotFoundError{
		ResourceType: "group",
		ID:           groupID,
	}
}

// UpdateGroup updates an existing group
func (c *Client) UpdateGroup(group *Group) (*Group, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_group":   group,
	})
	if err != nil {
		return nil, err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/update_group", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var result struct {
		Group Group `json:"group"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return &result.Group, nil
}

// DeleteGroup deletes a group by ID
func (c *Client) DeleteGroup(groupID string) error {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_group": map[string]string{
			"group_id": groupID,
		},
	})
	if err != nil {
		return err
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/delete_group", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// GetGroups retrieves all groups
func (c *Client) GetGroups() ([]Group, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_groups", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API request failed with status code: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	var result struct {
		Groups []Group `json:"groups"`
	}
	if err := json.Unmarshal(bodyBytes, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Groups, nil
}

// GetUsersGroup gets the users in a group
func (c *Client) GetUsersGroup(userID string) ([]GroupUser, error) {
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

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_users_group", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result GetUsersGroupResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Users, nil
}

// AddUserToGroup adds a user to a group
func (c *Client) AddUserToGroup(userID string, groupID string) error {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_group": map[string]interface{}{
			"group_id": groupID,
		},
		"target_user": map[string]interface{}{
			"user_id": userID,
		},
	})
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/add_user_group", "application/json", bytes.NewBuffer(body))
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

// RemoveUserFromGroup removes a user from a group
func (c *Client) RemoveUserFromGroup(userID string, groupID string) error {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_group": map[string]interface{}{
			"group_id": groupID,
		},
		"target_user": map[string]interface{}{
			"user_id": userID,
		},
	})
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/remove_user_group", "application/json", bytes.NewBuffer(body))
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

// GetGroupImages retrieves all images for a group
func (c *Client) GetGroupImages(groupID string) ([]GroupImage, error) {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_group": map[string]interface{}{
			"group_id": groupID,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/get_images_group", "application/json", bytes.NewBuffer(body))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Images []GroupImage `json:"images"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Images, nil
}

// AddGroupImage adds an image to a group
func (c *Client) AddGroupImage(groupID string, imageID string) error {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_group": map[string]interface{}{
			"group_id": groupID,
		},
		"target_image": map[string]interface{}{
			"image_id": imageID,
		},
	})
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/add_images_group", "application/json", bytes.NewBuffer(body))
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

// RemoveGroupImage removes an image from a group
func (c *Client) RemoveGroupImage(groupID string, imageID string) error {
	body, err := json.Marshal(map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_group": map[string]interface{}{
			"group_id": groupID,
		},
		"target_image": map[string]interface{}{
			"image_id": imageID,
		},
	})
	if err != nil {
		return fmt.Errorf("error marshaling request body: %v", err)
	}

	resp, err := c.HTTPClient.Post(c.BaseURL+"/api/public/remove_images_group", "application/json", bytes.NewBuffer(body))
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
