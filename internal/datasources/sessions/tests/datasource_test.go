package tests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/providerserver"
	"github.com/hashicorp/terraform-plugin-go/tfprotov6"
	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"terraform-provider-kasm/internal/client"
	"terraform-provider-kasm/internal/provider"
)

// Unit test for the sessions data source
func TestSessionsDataSource(t *testing.T) {
	// Create a mock server that returns a predefined response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := client.GetKasmsResponse{
			CurrentTime: "2023-01-01T00:00:00Z",
			Kasms: []client.Kasm{
				{
					ExpirationDate:      "2023-01-02T00:00:00Z",
					ContainerIP:         "10.0.0.1",
					StartDate:           "2023-01-01T00:00:00Z",
					Token:               "token123",
					ImageID:             "image123",
					ViewOnlyToken:       "viewtoken123",
					Cores:               2.0,
					Hostname:            "kasm-123",
					KasmID:              "kasm123",
					IsPersistentProfile: true,
					Memory:              1024,
					OperationalStatus:   "running",
					ContainerID:         "container123",
					Port:                8443,
					KeepaliveDate:       "2023-01-01T00:30:00Z",
					UserID:              "user123",
					ShareID:             "share123",
					Host:                "host123",
					ServerID:            "server123",
					PortMap: client.PortMap{
						Audio: struct {
							Port int    `json:"port"`
							Path string `json:"path"`
						}{
							Port: 8000,
							Path: "/audio",
						},
						VNC: struct {
							Port int    `json:"port"`
							Path string `json:"path"`
						}{
							Port: 5900,
							Path: "/vnc",
						},
						AudioInput: struct {
							Port int    `json:"port"`
							Path string `json:"path"`
						}{
							Port: 8001,
							Path: "/audioinput",
						},
					},
					Image: client.KasmImage{
						ImageID:      "image123",
						Name:         "Ubuntu",
						FriendlyName: "Ubuntu Desktop",
						ImageSrc:     "https://example.com/ubuntu.png",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create provider factories with the mock server URL
	providerFactories := map[string]func() (tfprotov6.ProviderServer, error){
		"kasm": providerserver.NewProtocol6WithError(provider.New(server.URL)),
	}

	// Run the test
	resource.Test(t, resource.TestCase{
		IsUnitTest:               true,
		ProtoV6ProviderFactories: providerFactories,
		Steps: []resource.TestStep{
			{
				Config: `
provider "kasm" {
  base_url = "` + server.URL + `"
  api_key = "test_key"
  api_secret = "test_secret"
}

data "kasm_sessions" "test" {}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kasm_sessions.test", "id", "kasm_sessions"),
					resource.TestCheckResourceAttr("data.kasm_sessions.test", "current_time", "2023-01-01T00:00:00Z"),
					resource.TestCheckResourceAttr("data.kasm_sessions.test", "sessions.#", "1"),
					resource.TestCheckResourceAttr("data.kasm_sessions.test", "sessions.0.kasm_id", "kasm123"),
					resource.TestCheckResourceAttr("data.kasm_sessions.test", "sessions.0.operational_status", "running"),
					resource.TestCheckResourceAttr("data.kasm_sessions.test", "sessions.0.image_name", "Ubuntu"),
					resource.TestCheckResourceAttr("data.kasm_sessions.test", "sessions.0.image_friendly_name", "Ubuntu Desktop"),
					resource.TestCheckResourceAttr("data.kasm_sessions.test", "sessions_map.kasm123.kasm_id", "kasm123"),
				),
			},
		},
	})
}
