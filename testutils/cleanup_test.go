//go:build acceptance
// +build acceptance

package testutils

import (
	"testing"
)

// TestCleanupResources is a utility test that can be run to clean up any resources
// created during testing, such as test images
func TestCleanupResources(t *testing.T) {
	// Clean up any test images that were created
	CleanupTestImage(t)

	// Clean up any existing sessions
	CleanupExistingSessions(t)

	t.Log("Cleanup completed successfully")
}
