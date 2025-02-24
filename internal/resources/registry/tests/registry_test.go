//go:build acceptance
// +build acceptance

package tests

import (
	"fmt"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/hashicorp/terraform-plugin-testing/plancheck"
	"terraform-provider-kasm/testutils"
)

const (
	// Default Kasm registry that supports channels/versions
	defaultKasmRegistryURL = "https://registry.kasmweb.com/"
	// Linuxserver registry for basic add/remove tests
	linuxserverRegistryURL = "https://kasmregistry.linuxserver.io/"
)

// Test for basic registry operations with channel support
func TestAccKasmRegistry_basic(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("%d", time.Now().Unix())
	registryName := fmt.Sprintf("test-registry-%s", uniqueIdentifier)
	resourceName := "kasm_registry.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmRegistryConfig_basic(registryName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "url", defaultKasmRegistryURL),
					resource.TestCheckResourceAttr(resourceName, "channel", "stable"),
				),
			},
		},
	})
}

// Configuration for lifecycle tests (using linuxserver registry)
func testAccKasmRegistryConfig_lifecycle(name string) string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_registry" "test" {
    url = "%s"
    channel = "stable"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"),
		os.Getenv("KASM_API_SECRET"), linuxserverRegistryURL)
}

func TestAccKasmRegistry_validation(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccKasmRegistryConfig_invalidURL(),
				ExpectError: regexp.MustCompile(`URL must start with http:// or https://`),
			},
		},
	})
}

// Test for basic addition and removal of a registry
func TestAccKasmRegistry_lifecycle(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("%d", time.Now().Unix())
	registryName := fmt.Sprintf("test-registry-%s", uniqueIdentifier)
	resourceName := "kasm_registry.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmRegistryConfig_lifecycle(registryName),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "url", linuxserverRegistryURL),
					resource.TestCheckResourceAttrSet(resourceName, "id"),
				),
			},
			{
				Config: testAccKasmRegistryConfig_empty(),
				ConfigPlanChecks: resource.ConfigPlanChecks{
					PreApply: []plancheck.PlanCheck{
						plancheck.ExpectResourceAction(resourceName, plancheck.ResourceActionDestroy),
					},
				},
			},
		},
	})
}

func TestAccKasmRegistry_dataSource(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmRegistryDataSourceConfig(),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttrSet("data.kasm_registries.all", "registries.#"),
					resource.TestCheckResourceAttrSet("data.kasm_registries.all", "registries.0.id"),
					resource.TestCheckResourceAttrSet("data.kasm_registries.all", "registries.0.url"),
					// Make channel check more flexible
					resource.TestMatchResourceAttr(
						"data.kasm_registries.all",
						"registries.0.channel",
						regexp.MustCompile(`^(stable|beta|develop|1\.\d+\.\d+)$`),
					),
				),
			},
		},
	})
}

func TestAccKasmRegistry_update(t *testing.T) {
	uniqueIdentifier := fmt.Sprintf("%d", time.Now().Unix())
	registryName := fmt.Sprintf("test-registry-%s", uniqueIdentifier)
	resourceName := "kasm_registry.test"

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: testAccKasmRegistryConfig_channel(registryName, "stable"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "channel", "stable"),
				),
			},
			{
				Config: testAccKasmRegistryConfig_channel(registryName, "beta"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "channel", "beta"),
				),
			},
			{
				Config: testAccKasmRegistryConfig_version(registryName, "1.16.0"),
				Check: resource.ComposeTestCheckFunc(
					resource.TestCheckResourceAttr(resourceName, "channel", "1.16.0"),
				),
			},
		},
	})
}

func TestAccKasmRegistry_errors(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testutils.TestAccPreCheck(t) },
		ProtoV6ProviderFactories: testutils.TestAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      testAccKasmRegistryConfig_missing_required(),
				ExpectError: regexp.MustCompile(`The argument "url" is required`),
			},
			{
				Config:      testAccKasmRegistryConfig_invalid_channel(),
				ExpectError: regexp.MustCompile(`Channel must be a valid Kasm version`),
			},
		},
	})
}

// Configuration for basic tests (using default Kasm registry)
func testAccKasmRegistryConfig_basic(name string) string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_registry" "test" {
    url = "%s"
    channel = "stable"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"),
		os.Getenv("KASM_API_SECRET"), defaultKasmRegistryURL)
}

func testAccKasmRegistryConfig_invalidURL() string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_registry" "test" {
    url = "invalid-url"
    channel = "stable"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}

func testAccKasmRegistryConfig_empty() string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}

func testAccKasmRegistryDataSourceConfig() string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

data "kasm_registries" "all" {}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}

func testAccKasmRegistryConfig_channel(name, channel string) string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_registry" "test" {
    url = "%s"
    channel = "%s"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"),
		os.Getenv("KASM_API_SECRET"), defaultKasmRegistryURL, channel)
}

// Add the missing test configuration functions
func testAccKasmRegistryConfig_invalid_channel() string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_registry" "test" {
    url = "https://registry.kasmweb.com/"
    channel = "invalid-channel"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}

func testAccKasmRegistryConfig_missing_required() string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_registry" "test" {
    channel = "stable"
    # url is required but missing
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"))
}

func testAccKasmRegistryConfig_version(name, version string) string {
	return fmt.Sprintf(`
provider "kasm" {
    base_url = "%s"
    api_key = "%s"
    api_secret = "%s"
    insecure = true
}

resource "kasm_registry" "test" {
    url = "https://registry.kasmweb.com/"
    channel = "%s"
}
`, os.Getenv("KASM_BASE_URL"), os.Getenv("KASM_API_KEY"), os.Getenv("KASM_API_SECRET"), version)
}
