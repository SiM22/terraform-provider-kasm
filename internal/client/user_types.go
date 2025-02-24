package client

// User represents a Kasm user
type User struct {
	UserID           string            `json:"user_id,omitempty"`
	Username         string            `json:"username"`
	FirstName        string            `json:"first_name"`
	LastName         string            `json:"last_name"`
	Password         string            `json:"password,omitempty"`
	Locked           bool              `json:"locked"`
	Disabled         bool              `json:"disabled"`
	Organization     string            `json:"organization,omitempty"`
	Phone            string            `json:"phone,omitempty"`
	Groups           []Group           `json:"groups,omitempty"`
	Attributes       map[string]string `json:"attributes,omitempty"`
	AuthorizedImages []string          `json:"authorized_images,omitempty"`
}

// UserAttributes represents user attributes
type UserAttributes struct {
	UserID            string            `json:"user_id,omitempty"`
	PreferredLanguage string            `json:"preferred_language,omitempty"`
	Theme             string            `json:"theme,omitempty"`
	DateFormat        string            `json:"date_format,omitempty"`
	TimeFormat        string            `json:"time_format,omitempty"`
	TimeZone          string            `json:"time_zone,omitempty"`
	Attributes        map[string]string `json:"attributes,omitempty"`
}

// GetUserAttributesResponse represents the response from the get_user_attributes API endpoint
type GetUserAttributesResponse struct {
	UserAttributes UserAttributes `json:"user_attributes"`
}

// UpdateUserAttributesRequest represents the request to update user attributes
type UpdateUserAttributesRequest struct {
	UserID     string            `json:"user_id"`
	Attributes map[string]string `json:"attributes"`
}

// UpdateUserAuthorizedImagesRequest represents the request to update user authorized images
type UpdateUserAuthorizedImagesRequest struct {
	UserID           string   `json:"user_id"`
	AuthorizedImages []string `json:"authorized_images"`
}

// GetUserAuthorizedImagesResponse represents the response from the get_user_authorized_images API endpoint
type GetUserAuthorizedImagesResponse struct {
	AuthorizedImages []string `json:"authorized_images"`
}
