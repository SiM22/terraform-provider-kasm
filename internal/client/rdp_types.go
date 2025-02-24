package client

// RDPConnectionType represents the type of RDP connection info to retrieve
type RDPConnectionType string

const (
	RDPConnectionTypeFile RDPConnectionType = "file"
	RDPConnectionTypeURL  RDPConnectionType = "url"
)

// RDPConnectionResponse represents the response from the get_rdp_client_connection_info API endpoint
type RDPConnectionResponse struct {
	File string `json:"file,omitempty"`
	URL  string `json:"url,omitempty"`
}
