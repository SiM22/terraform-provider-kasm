package client

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// GetRegistries retrieves all registries
func (c *Client) GetRegistries() ([]Registry, error) {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/get_registries", payload)
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
		Registries []Registry `json:"registries"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, string(body))
	}
	return result.Registries, nil
}

// CreateRegistryImage creates a new image from a registry workspace
func (c *Client) CreateRegistryImage(workspace RegistryWorkspace) (*Image, error) {
	image := ImageFromRegistryWorkspace(workspace)

	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_image":   image,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/create_image", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Image *Image `json:"image"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Image, nil
}

// GetRegistryImage retrieves an image from the registry
func (c *Client) GetRegistryImage(imageID string) (*Image, error) {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_image": map[string]string{
			"image_id": imageID,
		},
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/get_image", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Image *Image `json:"image"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Image, nil
}

// UpdateRegistryImage updates an existing registry image
func (c *Client) UpdateRegistryImage(image *Image) (*Image, error) {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_image":   image,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/update_image", payload)
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Image *Image `json:"image"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return result.Image, nil
}

// DeleteRegistryImage deletes an image from the registry
func (c *Client) DeleteRegistryImage(imageID string) error {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_image": map[string]string{
			"image_id": imageID,
		},
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/delete_image", payload)
	if err != nil {
		return fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	return nil
}

// ListRegistryImages retrieves all registry images
func (c *Client) ListRegistryImages() ([]RegistryImage, error) {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/get_images", payload)
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
		Images []RegistryImage `json:"images"`
	}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("error decoding response: %v, body: %s", err, string(body))
	}
	return result.Images, nil
}

// CreateRegistry creates a new registry
func (c *Client) CreateRegistry(request *CreateRegistryRequest) error {
	payload := map[string]interface{}{
		"api_key":         c.APIKey,
		"api_key_secret":  c.APISecret,
		"registry":        request.Registry,
		"override_schema": request.OverrideSchema,
		"channel":         request.Channel,
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/create_registry", payload)
	if err != nil {
		return fmt.Errorf("error creating registry: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to create registry, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// DeleteRegistry deletes a registry
func (c *Client) DeleteRegistry(registryID string) error {
	payload := map[string]interface{}{
		"api_key":        c.APIKey,
		"api_key_secret": c.APISecret,
		"target_registry": map[string]string{
			"registry_id": registryID,
		},
	}

	resp, err := c.doRequestLegacy("POST", "/api/public/delete_registry", payload)
	if err != nil {
		return fmt.Errorf("error deleting registry: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to delete registry, status: %d, response: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// ImageFromRegistryWorkspace converts a RegistryWorkspace to an Image
func ImageFromRegistryWorkspace(workspace RegistryWorkspace) *Image {
	return &Image{
		Name:         workspace.Name,
		FriendlyName: workspace.Name,
		Description:  workspace.Description,
		ImageSrc:     workspace.ImageURL,
	}
}
