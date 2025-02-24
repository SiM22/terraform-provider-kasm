package client

// StagingConfig represents a Kasm staging configuration
type StagingConfig struct {
	StagingConfigID        string  `json:"staging_config_id"`
	ZoneID                 string  `json:"zone_id"`
	ZoneName               string  `json:"zone_name"`
	ImageID                string  `json:"image_id"`
	ImageFriendlyName      string  `json:"image_friendly_name"`
	NumSessions            int     `json:"num_sessions"`
	NumCurrentSessions     int     `json:"num_current_sessions"`
	Expiration             float64 `json:"expiration"`
	AllowKasmAudio         bool    `json:"allow_kasm_audio"`
	AllowKasmUploads       bool    `json:"allow_kasm_uploads"`
	AllowKasmDownloads     bool    `json:"allow_kasm_downloads"`
	AllowKasmClipboardDown bool    `json:"allow_kasm_clipboard_down"`
	AllowKasmClipboardUp   bool    `json:"allow_kasm_clipboard_up"`
	AllowKasmMicrophone    bool    `json:"allow_kasm_microphone"`
}

// CreateStagingConfigRequest represents the request to create a staging config
type CreateStagingConfigRequest struct {
	ZoneID                 string  `json:"zone_id"`
	ImageID                string  `json:"image_id"`
	NumSessions            int     `json:"num_sessions"`
	Expiration             float64 `json:"expiration"`
	AllowKasmAudio         bool    `json:"allow_kasm_audio"`
	AllowKasmUploads       bool    `json:"allow_kasm_uploads"`
	AllowKasmDownloads     bool    `json:"allow_kasm_downloads"`
	AllowKasmClipboardDown bool    `json:"allow_kasm_clipboard_down"`
	AllowKasmClipboardUp   bool    `json:"allow_kasm_clipboard_up"`
	AllowKasmMicrophone    bool    `json:"allow_kasm_microphone"`
}

// UpdateStagingConfigRequest represents the request to update a staging config
type UpdateStagingConfigRequest struct {
	StagingConfigID        string   `json:"staging_config_id"`
	ZoneID                 string   `json:"zone_id,omitempty"`
	ImageID                string   `json:"image_id,omitempty"`
	NumSessions            *int     `json:"num_sessions,omitempty"`
	Expiration             *float64 `json:"expiration,omitempty"`
	AllowKasmAudio         *bool    `json:"allow_kasm_audio,omitempty"`
	AllowKasmUploads       *bool    `json:"allow_kasm_uploads,omitempty"`
	AllowKasmDownloads     *bool    `json:"allow_kasm_downloads,omitempty"`
	AllowKasmClipboardDown *bool    `json:"allow_kasm_clipboard_down,omitempty"`
	AllowKasmClipboardUp   *bool    `json:"allow_kasm_clipboard_up,omitempty"`
	AllowKasmMicrophone    *bool    `json:"allow_kasm_microphone,omitempty"`
}
