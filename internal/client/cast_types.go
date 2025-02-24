package client

// CastConfig represents a session casting configuration
type CastConfig struct {
	ID                   string   `json:"cast_config_id,omitempty"`
	CastingConfigName    string   `json:"casting_config_name"`
	ImageID              string   `json:"image_id"`
	ImageFriendlyName    string   `json:"image_friendly_name,omitempty"`
	AllowedReferrers     []string `json:"allowed_referrers,omitempty"`
	LimitSessions        bool     `json:"limit_sessions"`
	SessionRemaining     int      `json:"session_remaining"`
	LimitIPs            bool     `json:"limit_ips"`
	IPRequestLimit      int      `json:"ip_request_limit"`
	IPRequestSeconds    int      `json:"ip_request_seconds"`
	ErrorURL            string   `json:"error_url,omitempty"`
	EnableSharing       bool     `json:"enable_sharing"`
	DisableControlPanel bool     `json:"disable_control_panel"`
	DisableTips         bool     `json:"disable_tips"`
	DisableFixedRes     bool     `json:"disable_fixed_res"`
	Key                 string   `json:"key"`
	AllowAnonymous      bool     `json:"allow_anonymous"`
	GroupID             string   `json:"group_id,omitempty"`
	RequireRecaptcha    bool     `json:"require_recaptcha"`
	KasmURL            string   `json:"kasm_url,omitempty"`
	DynamicKasmURL     bool     `json:"dynamic_kasm_url"`
	DynamicDockerNetwork bool   `json:"dynamic_docker_network"`
	AllowResume         bool    `json:"allow_resume"`
	EnforceClientSettings bool  `json:"enforce_client_settings"`
	AllowKasmAudio      bool    `json:"allow_kasm_audio"`
	AllowKasmUploads    bool    `json:"allow_kasm_uploads"`
	AllowKasmDownloads  bool    `json:"allow_kasm_downloads"`
	AllowClipboardDown  bool    `json:"allow_kasm_clipboard_down"`
	AllowClipboardUp    bool    `json:"allow_kasm_clipboard_up"`
	AllowMicrophone     bool    `json:"allow_kasm_microphone"`
	ValidUntil          string  `json:"valid_until,omitempty"`
	AllowSharing        bool    `json:"allow_kasm_sharing"`
	AudioDefaultOn      bool    `json:"kasm_audio_default_on"`
	IMEModeDefaultOn    bool    `json:"kasm_ime_mode_default_on"`
}

type createCastConfigRequest struct {
	APIKey           string     `json:"api_key"`
	APIKeySecret     string     `json:"api_key_secret"`
	TargetCastConfig CastConfig `json:"target_cast_config"`
}

type createCastConfigResponse struct {
	CastConfig CastConfig `json:"cast_config"`
}

type getCastConfigRequest struct {
	APIKey       string `json:"api_key"`
	APIKeySecret string `json:"api_key_secret"`
	CastConfigID string `json:"cast_config_id"`
}

type getCastConfigResponse struct {
	CastConfig CastConfig `json:"cast_config"`
}

type updateCastConfigRequest struct {
	APIKey           string     `json:"api_key"`
	APIKeySecret     string     `json:"api_key_secret"`
	TargetCastConfig CastConfig `json:"target_cast_config"`
}

type updateCastConfigResponse struct {
	CastConfig CastConfig `json:"cast_config"`
}

type deleteCastConfigRequest struct {
	APIKey       string `json:"api_key"`
	APIKeySecret string `json:"api_key_secret"`
	CastConfigID string `json:"cast_config_id"`
}
