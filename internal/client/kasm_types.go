package client

// Kasm represents a Kasm session
type Kasm struct {
	ExpirationDate        string         `json:"expiration_date"`
	ContainerIP           string         `json:"container_ip"`
	StartDate             string         `json:"start_date"`
	PointOfPresence       interface{}    `json:"point_of_presence"`
	Token                 string         `json:"token"`
	ImageID               string         `json:"image_id"`
	ViewOnlyToken         string         `json:"view_only_token"`
	Cores                 float64        `json:"cores"`
	Hostname              string         `json:"hostname"`
	KasmID                string         `json:"kasm_id"`
	PortMap               PortMap        `json:"port_map"`
	Image                 KasmImage      `json:"image"`
	IsPersistentProfile   bool           `json:"is_persistent_profile"`
	Memory                int64          `json:"memory"`
	OperationalStatus     string         `json:"operational_status"`
	ClientSettings        ClientSettings `json:"client_settings"`
	ContainerID           string         `json:"container_id"`
	Port                  int            `json:"port"`
	KeepaliveDate         string         `json:"keepalive_date"`
	UserID                string         `json:"user_id"`
	PersistentProfileMode interface{}    `json:"persistent_profile_mode"`
	ShareID               string         `json:"share_id"`
	Host                  string         `json:"host"`
	ServerID              string         `json:"server_id"`
}

// KasmImage represents the image information for a Kasm session
type KasmImage struct {
	ImageID      string `json:"image_id"`
	Name         string `json:"name"`
	ImageSrc     string `json:"image_src"`
	FriendlyName string `json:"friendly_name"`
}

// PortMap represents the port mapping configuration for a Kasm session
type PortMap struct {
	Audio struct {
		Port int    `json:"port"`
		Path string `json:"path"`
	} `json:"audio"`
	VNC struct {
		Port int    `json:"port"`
		Path string `json:"path"`
	} `json:"vnc"`
	AudioInput struct {
		Port int    `json:"port"`
		Path string `json:"path"`
	} `json:"audio_input"`
	Uploads struct {
		Port int    `json:"port"`
		Path string `json:"path"`
	} `json:"uploads"`
}

// ClientSettings represents the client configuration for a Kasm session
type ClientSettings struct {
	AllowKasmAudio             bool    `json:"allow_kasm_audio"`
	IdleDisconnect             float64 `json:"idle_disconnect"` // Changed from int to float64
	LockSharingVideoMode       bool    `json:"lock_sharing_video_mode"`
	AllowPersistentProfile     bool    `json:"allow_persistent_profile"`
	AllowKasmClipboardDown     bool    `json:"allow_kasm_clipboard_down"`
	AllowKasmMicrophone        bool    `json:"allow_kasm_microphone"`
	AllowKasmDownloads         bool    `json:"allow_kasm_downloads"`
	KasmAudioDefaultOn         bool    `json:"kasm_audio_default_on"`
	AllowPointOfPresence       bool    `json:"allow_point_of_presence"`
	AllowKasmUploads           bool    `json:"allow_kasm_uploads"`
	AllowKasmClipboardUp       bool    `json:"allow_kasm_clipboard_up"`
	EnableWebp                 bool    `json:"enable_webp"`
	AllowKasmSharing           bool    `json:"allow_kasm_sharing"`
	AllowKasmClipboardSeamless bool    `json:"allow_kasm_clipboard_seamless"`
}

// KasmStatusResponse represents the response from the get_kasm_status API endpoint
type KasmStatusResponse struct {
	CurrentTime         string `json:"current_time,omitempty"`
	KasmURL             string `json:"kasm_url,omitempty"`
	Kasm                *Kasm  `json:"kasm,omitempty"`
	ErrorMessage        string `json:"error_message,omitempty"`
	OperationalMessage  string `json:"operational_message,omitempty"`
	OperationalProgress int    `json:"operational_progress,omitempty"`
	OperationalStatus   string `json:"operational_status,omitempty"`
}

// CreateKasmResponse represents the response from the request_kasm API endpoint
type CreateKasmResponse struct {
	KasmID       string `json:"kasm_id"`
	Status       string `json:"status"`
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	SessionToken string `json:"session_token"`
	KasmURL      string `json:"kasm_url"`
	ShareID      string `json:"share_id"`
	ErrorMessage string `json:"error_message,omitempty"`
}

// JoinKasmResponse represents the response from the join_kasm API endpoint
type JoinKasmResponse struct {
	CurrentTime  string `json:"current_time"`
	SessionToken string `json:"session_token"`
	UserID       string `json:"user_id"`
	ErrorMessage string `json:"error_message,omitempty"`
	Kasm         struct {
		PortMap struct {
			VNC struct {
				Port int    `json:"port"`
				Path string `json:"path"`
			} `json:"vnc"`
			Audio struct {
				Port int    `json:"port"`
				Path string `json:"path"`
			} `json:"audio"`
		} `json:"port_map"`
		Port     int    `json:"port"`
		Hostname string `json:"hostname"`
		Image    struct {
			ImageID      string `json:"image_id"`
			Name         string `json:"name"`
			ImageSrc     string `json:"image_src"`
			FriendlyName string `json:"friendly_name"`
		} `json:"image"`
		ViewOnlyToken string `json:"view_only_token"`
		User          struct {
			Username string `json:"username"`
		} `json:"user"`
		ShareID        string         `json:"share_id"`
		Host           string         `json:"host"`
		ClientSettings ClientSettings `json:"client_settings"`
		KasmID         string         `json:"kasm_id"`
	} `json:"kasm"`
	Username string `json:"username"`
	KasmURL  string `json:"kasm_url"`
}

// GetKasmsResponse represents the response from the get_kasms API endpoint
type GetKasmsResponse struct {
	CurrentTime string `json:"current_time"`
	Kasms       []Kasm `json:"kasms"`
}

// KasmServer represents server information for a Kasm session
type KasmServer struct {
	Port     int    `json:"port"`
	Hostname string `json:"hostname"`
	ZoneName string `json:"zone_name"`
	Provider string `json:"provider"`
}
