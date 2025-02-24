package client

// CastingConfig represents a Kasm cast configuration
type CastingConfig struct {
	CastingConfigName      string   `json:"casting_config_name"`
	CastConfigID           string   `json:"cast_config_id"`
	ImageID                string   `json:"image_id"`
	ImageFriendlyName      string   `json:"image_friendly_name"`
	AllowedReferrers       []string `json:"allowed_referrers"`
	LimitSessions          bool     `json:"limit_sessions"`
	SessionRemaining       int      `json:"session_remaining"`
	LimitIPs               bool     `json:"limit_ips"`
	IPRequestLimit         int      `json:"ip_request_limit"`
	IPRequestSeconds       int      `json:"ip_request_seconds"`
	ErrorURL               *string  `json:"error_url"`
	EnableSharing          bool     `json:"enable_sharing"`
	DisableControlPanel    bool     `json:"disable_control_panel"`
	DisableTips            bool     `json:"disable_tips"`
	DisableFixedRes        bool     `json:"disable_fixed_res"`
	Key                    string   `json:"key"`
	AllowAnonymous         bool     `json:"allow_anonymous"`
	GroupID                *string  `json:"group_id"`
	RequireRecaptcha       bool     `json:"require_recaptcha"`
	GroupName              *string  `json:"group_name"`
	KasmURL                *string  `json:"kasm_url"`
	DynamicKasmURL         bool     `json:"dynamic_kasm_url"`
	DynamicDockerNetwork   bool     `json:"dynamic_docker_network"`
	AllowResume            bool     `json:"allow_resume"`
	EnforceClientSettings  bool     `json:"enforce_client_settings"`
	AllowKasmAudio         bool     `json:"allow_kasm_audio"`
	AllowKasmUploads       bool     `json:"allow_kasm_uploads"`
	AllowKasmDownloads     bool     `json:"allow_kasm_downloads"`
	AllowKasmClipboardDown bool     `json:"allow_kasm_clipboard_down"`
	AllowKasmClipboardUp   bool     `json:"allow_kasm_clipboard_up"`
	AllowKasmMicrophone    bool     `json:"allow_kasm_microphone"`
	ValidUntil             string   `json:"valid_until"`
	AllowKasmSharing       bool     `json:"allow_kasm_sharing"`
	KasmAudioDefaultOn     bool     `json:"kasm_audio_default_on"`
	KasmIMEModeDefaultOn   bool     `json:"kasm_ime_mode_default_on"`
}

// CreateCastingConfigRequest represents the request to create a cast config
type CreateCastingConfigRequest struct {
	TargetCastConfig CastingConfig `json:"target_cast_config"`
}

// UpdateCastingConfigRequest represents the request to update a cast config
type UpdateCastingConfigRequest struct {
	TargetCastConfig CastingConfig `json:"target_cast_config"`
}

// DeleteCastingConfigRequest represents the request to delete a cast config
type DeleteCastingConfigRequest struct {
	CastConfigID      string `json:"cast_config_id"`
	CastingConfigName string `json:"casting_config_name"`
}
