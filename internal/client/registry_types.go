package client

// Registry represents a Kasm registry
type Registry struct {
	RegistryID    string                      `json:"registry_id"`
	RegistryURL   string                      `json:"registry_url"`
	DoAutoUpdate  bool                        `json:"do_auto_update"`
	SchemaVersion string                      `json:"schema_version"`
	IsVerified    bool                        `json:"is_verified"`
	Channel       string                      `json:"channel"`
	Workspaces    []RegistryWorkspaceResponse `json:"workspaces"`
}

// RegistryWorkspace represents a workspace in a registry
type RegistryWorkspace struct {
	Name        string `json:"name"`
	ImageURL    string `json:"image_url"`
	Description string `json:"description"`
}

// RegistryWorkspaceResponse represents a workspace in a registry response
type RegistryWorkspaceResponse struct {
	Compatibility  []Compatibility        `json:"compatibility"`
	DockerRegistry string                 `json:"docker_registry"`
	FriendlyName   string                 `json:"friendly_name"`
	Description    string                 `json:"description"`
	RunConfig      map[string]interface{} `json:"run_config"`
	ExecConfig     map[string]interface{} `json:"exec_config"`
	ImageSrc       string                 `json:"image_src"`
}

// Compatibility represents the compatibility information for a workspace image
type Compatibility struct {
	Image              string   `json:"image"`
	Version            string   `json:"version"`
	AvailableTags      []string `json:"available_tags"`
	UncompressedSizeMB int      `json:"uncompressed_size_mb"`
}

// CreateRegistryRequest represents the request to create a registry
type CreateRegistryRequest struct {
	Registry       string `json:"registry"`
	OverrideSchema string `json:"override_schema,omitempty"`
	Channel        string `json:"channel"`
}
