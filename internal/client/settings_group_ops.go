package client

import (
	"fmt"
)

type Setting struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

// SetSettingsGroup configures the settings group with the specified settings
func (c *Client) SetSettingsGroup(settings []Setting) error {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"settings":       settings,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/set_settings_group", payload)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}

// ConfigureDefaultSharingSettings sets up the default sharing configuration
func (c *Client) ConfigureDefaultSharingSettings(groupID string) error {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_group": map[string]string{
			"group_id": groupID,
		},
		"target_setting": map[string]interface{}{
			"name":  "allow_kasm_sharing",
			"value": "True",
		},
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/add_settings_group", payload)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	return nil
}
