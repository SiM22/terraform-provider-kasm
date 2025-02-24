package client

// Image represents a Kasm image as returned by the documented API
type Image struct {
	ImageID             string                 `json:"image_id,omitempty"`
	Name                string                 `json:"name"`
	FriendlyName        string                 `json:"friendly_name"`
	Description         string                 `json:"description"`
	Categories          []string               `json:"categories"`
	Memory              int64                  `json:"memory"`
	Cores               float64                `json:"cores"`
	CPUAllocationMethod string                 `json:"cpu_allocation_method"`
	DockerRegistry      string                 `json:"docker_registry"`
	DockerUser          string                 `json:"docker_user,omitempty"`
	DockerPassword      string                 `json:"docker_password,omitempty"`
	UncompressedSizeMB  int                    `json:"uncompressed_size_mb"`
	ImageType           string                 `json:"image_type"`
	Enabled             bool                   `json:"enabled"`
	Available           bool                   `json:"available"`
	ImageSrc            string                 `json:"image_src"`
	RunConfig           map[string]interface{} `json:"run_config"`
	ExecConfig          map[string]interface{} `json:"exec_config"`
	VolumeMappings      map[string]interface{} `json:"volume_mappings"`
	RestrictToNetwork   bool                   `json:"restrict_to_network"`
	RestrictToServer    bool                   `json:"restrict_to_server"`
	RestrictToZone      bool                   `json:"restrict_to_zone"`
	ServerID            string                 `json:"server_id,omitempty"`
	ZoneID              string                 `json:"zone_id,omitempty"`
	NetworkName         string                 `json:"network_name,omitempty"`
}

// RegistryImage represents a Kasm registry image from the get_images API endpoint
type RegistryImage struct {
	ImageID        string  `json:"image_id"`
	Name           string  `json:"name"`
	FriendlyName   string  `json:"friendly_name"`
	Description    string  `json:"description"`
	Memory         int64   `json:"memory"`
	Cores          float64 `json:"cores"`
	DockerRegistry string  `json:"docker_registry"`
}

// SessionRecording represents a Kasm session recording
type SessionRecording struct {
	RecordingID                 string                 `json:"recording_id"`
	AccountID                   string                 `json:"account_id"`
	SessionRecordingURL         string                 `json:"session_recording_url"`
	SessionRecordingMetadata    map[string]interface{} `json:"session_recording_metadata"`
	SessionRecordingDownloadURL string                 `json:"session_recording_download_url,omitempty"`
}

// CreateImageRequest represents the request to create a workspace image
type CreateImageRequest struct {
	ImageSrc           string  `json:"image_src"`
	Categories         string  `json:"categories"`
	RunConfig          string  `json:"run_config"`
	Description        string  `json:"description"`
	ExecConfig         string  `json:"exec_config"`
	FriendlyName       string  `json:"friendly_name"`
	DockerRegistry     string  `json:"docker_registry"`
	Name               string  `json:"name"`
	UncompressedSizeMB int     `json:"uncompressed_size_mb"`
	ImageType          string  `json:"image_type"`
	Enabled            bool    `json:"enabled"`
	Memory             int64   `json:"memory"`
	Cores              float64 `json:"cores"`
	GPUCount           int     `json:"gpu_count"`
	RequireGPU         *bool   `json:"require_gpu"`
	VolumeMappings     string  `json:"volume_mappings"`
}

// Request wrapper for API
type createImageAPIRequest struct {
	APIKey       string             `json:"api_key"`
	APIKeySecret string             `json:"api_key_secret"`
	TargetImage  CreateImageRequest `json:"target_image"`
}
