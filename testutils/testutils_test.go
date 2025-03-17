//go:build unit
// +build unit

package testutils

import (
	"encoding/json"
	"testing"
)

func TestVolumeMapppingsFormat(t *testing.T) {
	// Create the volume mappings structure
	volumeMappings := map[string]interface{}{
		"uploads": map[string]interface{}{
			"host_path":      "/home/kasm-user/uploads",
			"container_path": "/home/kasm-user/uploads",
			"mode":           "rw",
			"bind":           true,
			"uid":            1000,
			"gid":            1000,
		},
	}

	// Marshal to JSON
	volumeMappingsJSON, err := json.Marshal(volumeMappings)
	if err != nil {
		t.Fatalf("Error marshaling volume mappings: %v", err)
	}

	// Verify the JSON structure
	expected := `{"uploads":{"bind":true,"container_path":"/home/kasm-user/uploads","gid":1000,"host_path":"/home/kasm-user/uploads","mode":"rw","uid":1000}}`
	if string(volumeMappingsJSON) != expected {
		t.Errorf("Volume mappings JSON doesn't match expected format.\nGot: %s\nExpected: %s", string(volumeMappingsJSON), expected)
	}
}
