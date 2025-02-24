package client

// SessionToken represents a Kasm session token
type SessionToken struct {
	SessionToken     string `json:"session_token"`
	SessionTokenDate string `json:"session_token_date"`
	ExpiresAt        string `json:"expires_at"`
	SessionJWT       string `json:"session_jwt,omitempty"`
}

// CreateSessionTokenRequest represents the request to create a session token
type CreateSessionTokenRequest struct {
	TargetUser struct {
		UserID string `json:"user_id"`
	} `json:"target_user"`
}

// GetSessionTokenRequest represents the request to get a session token
type GetSessionTokenRequest struct {
	TargetSessionToken struct {
		SessionToken string `json:"session_token"`
	} `json:"target_session_token"`
}

// GetSessionTokensRequest represents the request to get all session tokens for a user
type GetSessionTokensRequest struct {
	TargetUser struct {
		UserID string `json:"user_id"`
	} `json:"target_user"`
}

// UpdateSessionTokenRequest represents the request to update a session token
type UpdateSessionTokenRequest struct {
	TargetSessionToken struct {
		SessionToken string `json:"session_token"`
	} `json:"target_session_token"`
}

// DeleteSessionTokenRequest represents the request to delete a session token
type DeleteSessionTokenRequest struct {
	TargetSessionToken struct {
		SessionToken string `json:"session_token"`
	} `json:"target_session_token"`
}

// DeleteSessionTokensRequest represents the request to delete all session tokens for a user
type DeleteSessionTokensRequest struct {
	TargetUser struct {
		UserID string `json:"user_id"`
	} `json:"target_user"`
}
