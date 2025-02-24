package client

// LicenseFeatures represents the features enabled by a license
type LicenseFeatures struct {
	AutoScaling       bool `json:"auto_scaling"`
	Branding          bool `json:"branding"`
	SessionStaging    bool `json:"session_staging"`
	SessionCasting    bool `json:"session_casting"`
	LogForwarding     bool `json:"log_forwarding"`
	DeveloperAPI      bool `json:"developer_api"`
	InjectSSHKeys     bool `json:"inject_ssh_keys"`
	SAML              bool `json:"saml"`
	LDAP              bool `json:"ldap"`
	SessionSharing    bool `json:"session_sharing"`
	LoginBanner       bool `json:"login_banner"`
	URLCategorization bool `json:"url_categorization"`
	UsageLimit        bool `json:"usage_limit"`
}

// License represents a Kasm license
type License struct {
	LicenseID   string          `json:"license_id"`
	Expiration  string          `json:"expiration"`
	IssuedAt    string          `json:"issued_at"`
	IssuedTo    string          `json:"issued_to"`
	Limit       int             `json:"limit"`
	IsVerified  bool            `json:"is_verified"`
	LicenseType string          `json:"license_type"`
	Features    LicenseFeatures `json:"features"`
	SKU         string          `json:"sku"`
}

// ActivateRequest represents the request to activate a license
type ActivateRequest struct {
	ActivationKey string `json:"activation_key"`
	Seats         *int   `json:"seats,omitempty"`
	IssuedTo      string `json:"issued_to,omitempty"`
}

// ActivateResponse represents the response from the activate API endpoint
type ActivateResponse struct {
	License License `json:"license"`
}

// GetLicensesResponse represents the response from the get_licenses API endpoint
type GetLicensesResponse struct {
	Success bool      `json:"success"`
	Code    int       `json:"code"`
	Message string    `json:"message"`
	Data    []License `json:"data"`
}
