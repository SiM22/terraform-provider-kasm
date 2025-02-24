package client

// ExecConfig represents the configuration for executing a command in a Kasm session
type ExecConfig struct {
	Cmd         string            `json:"cmd"`
	Environment map[string]string `json:"environment,omitempty"`
	Workdir     string            `json:"workdir,omitempty"`
	Privileged  bool              `json:"privileged,omitempty"`
	User        string            `json:"user,omitempty"`
}

// ExecCommandResponse represents the response from executing a command in a Kasm session
type ExecCommandResponse struct {
	Kasm struct {
		ExitCode int    `json:"exit_code"`
		Stdout   string `json:"stdout"`
		Stderr   string `json:"stderr"`
	} `json:"kasm"`
}
