package client

// Zone represents a Kasm deployment zone
type Zone struct {
	ZoneID             string `json:"zone_id"`
	ZoneName           string `json:"zone_name"`
	AutoScalingEnabled bool   `json:"auto_scaling_enabled"`
	AWSEnabled         bool   `json:"aws_enabled"`
	AWSRegion          string `json:"aws_region"`
	AWSAccessKeyID     string `json:"aws_access_key_id"`
	AWSSecretAccessKey string `json:"aws_secret_access_key"`
	EC2AgentAMIID      string `json:"ec2_agent_ami_id"`
}

// GetZonesRequest represents the request to get zones
type GetZonesRequest struct {
	Brief bool `json:"brief,omitempty"`
}
