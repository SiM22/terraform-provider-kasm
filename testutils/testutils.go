package testutils

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/joho/godotenv"
	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/internal/provider"
)

// LoadEnvFile loads environment variables from .env file
func LoadEnvFile(t *testing.T) {
	// Try to find .env file in current directory and parent directories
	dir, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	for {
		envFile := filepath.Join(dir, ".env")
		if _, err := os.Stat(envFile); err == nil {
			err = godotenv.Load(envFile)
			if err != nil {
				t.Fatalf("Error loading .env file: %v", err)
			}
			return
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	t.Log("No .env file found, using existing environment variables")
}

// TestAccPreCheck validates required environment variables exist
func TestAccPreCheck(t *testing.T) {
	LoadEnvFile(t)

	if v := os.Getenv("KASM_BASE_URL"); v == "" {
		t.Fatal("KASM_BASE_URL must be set for acceptance tests")
	}
	if v := os.Getenv("KASM_API_KEY"); v == "" {
		t.Fatal("KASM_API_KEY must be set for acceptance tests")
	}
	if v := os.Getenv("KASM_API_SECRET"); v == "" {
		t.Fatal("KASM_API_SECRET must be set for acceptance tests")
	}
}

// TestAccProtoV6ProviderFactories are used to instantiate a provider during
// acceptance tests
var TestAccProtoV6ProviderFactories = map[string]func() (tfprotov6.ProviderServer, error){
	"kasm": providerserver.NewProtocol6WithError(provider.New()),
}

// GetTestClient returns a configured Kasm client for testing
func GetTestClient(t *testing.T) *client.Client {
	return client.NewClient(
		os.Getenv("KASM_BASE_URL"),
		os.Getenv("KASM_API_KEY"),
		os.Getenv("KASM_API_SECRET"),
		true,
	)
}

// providerConfig (lowercase) for internal use
func providerConfig() string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}

// ProviderConfig (uppercase) for external use
func ProviderConfig() string {
	return providerConfig()
}

// User test configurations
func TestAccUserConfig(username string) string {
	return fmt.Sprintf(`
%s

resource "kasm_user" "test" {
    username     = "%s"
    password     = "TestImport123!"
    first_name   = "Test"
    last_name    = "Import"
    organization = "Import Testing"
    phone        = "1234567890"
    groups       = []
}`, ProviderConfig(), username)
}

// Group test configurations
func TestAccGroupConfig(name string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group" "test" {
    name        = "%s"
    priority    = 50
    description = "Test Import Group"
}`, ProviderConfig(), name)
}

// Combined user and group test configurations
func TestAccUserGroupConfig(username, groupName string) string {
	return fmt.Sprintf(`
%s

resource "kasm_group" "test" {
    name        = "%s"
    priority    = 50
    description = "Test group for user import"
}

resource "kasm_user" "test" {
    username     = "%s"
    password     = "TestImport123!"
    first_name   = "Test"
    last_name    = "Import"
    organization = "Import Testing"
    groups       = [kasm_group.test.name]
}`, ProviderConfig(), groupName, username)
}

// Empty config for import testing
func TestAccEmptyConfig() string {
	return ProviderConfig()
}

// TestCheckResourceExists verifies that a resource exists in state by checking an attribute
func TestCheckResourceExists(resourceName string) resource.TestCheckFunc {
	return resource.ComposeTestCheckFunc(
		resource.TestCheckResourceAttrSet(resourceName, "id"), // Change "id" to an actual attribute name to confirm existence
	)
}
