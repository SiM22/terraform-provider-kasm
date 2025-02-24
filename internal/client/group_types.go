package client

// Group represents a Kasm group
type Group struct {
	GroupID     string   `json:"group_id,omitempty"`
	Name        string   `json:"name"`
	Description string   `json:"description,omitempty"`
	Priority    int      `json:"priority,omitempty"`
	IsSystem    bool     `json:"is_system,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

// GroupUser represents a user in a group
type GroupUser struct {
	UserID       string `json:"user_id"`
	Username     string `json:"username"`
	FirstName    string `json:"first_name"`
	LastName     string `json:"last_name"`
	Organization string `json:"organization"`
}

// GetUsersGroupResponse represents the response from the get_users_group API endpoint
type GetUsersGroupResponse struct {
	Users []GroupUser `json:"users"`
}

// GroupImage represents an image authorized for a group
type GroupImage struct {
	GroupImageID      string `json:"group_image_id"`
	GroupID           string `json:"group_id"`
	ImageID           string `json:"image_id"`
	ImageName         string `json:"image_name"`
	GroupName         string `json:"group_name"`
	ImageFriendlyName string `json:"image_friendly_name"`
	ImageSrc          string `json:"image_src"`
}
