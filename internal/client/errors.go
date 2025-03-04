package client

import "fmt"

// NotFoundError represents a resource not found error
type NotFoundError struct {
	ResourceType string
	ID           string
}

func (e *NotFoundError) Error() string {
	return fmt.Sprintf("%s not found: %s", e.ResourceType, e.ID)
}

// IsGroupNotFoundError checks if the error is due to a group not being found
func IsGroupNotFoundError(err error) bool {
	if err == nil {
		return false
	}
	if nfe, ok := err.(*NotFoundError); ok {
		return nfe.ResourceType == "group"
	}
	return false
}

// APIError represents an error returned by the API
type APIError struct {
	StatusCode   int
	ErrorMessage string
	Response     string
}

func (e *APIError) Error() string {
	if e.ErrorMessage != "" {
		return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.ErrorMessage)
	}
	return fmt.Sprintf("API error (status %d): %s", e.StatusCode, e.Response)
}

// UnauthorizedError represents an unauthorized access error
type UnauthorizedError struct {
	Message string
}

func (e *UnauthorizedError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("Unauthorized: %s", e.Message)
	}
	return "Unauthorized"
}

// IsResourceUnavailableError checks if the error is due to resources being unavailable
func IsResourceUnavailableError(err error) bool {
	if err == nil {
		return false
	}
	if apiErr, ok := err.(*APIError); ok {
		return apiErr.ErrorMessage == "No resources are available to create the requested Kasm. Please try again later or contact an Administrator"
	}
	return false
}
