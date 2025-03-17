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

// Unit test for the session_status data source
func TestSessionStatusDataSource(t *testing.T) {
	// Create a mock server that returns a predefined response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := client.KasmStatusResponse{
			CurrentTime:         "2023-01-01T00:00:00Z",
			KasmURL:             "https://example.com/kasm/123",
			OperationalStatus:   "running",
			OperationalMessage:  "Session is running",
			OperationalProgress: 100,
			Kasm: &client.Kasm{
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
				Image: client.KasmImage{
					ImageID:      "image123",
					Name:         "Ubuntu",
					FriendlyName: "Ubuntu Desktop",
					ImageSrc:     "https://example.com/ubuntu.png",
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

data "kasm_session_status" "test" {
  kasm_id = "kasm123"
  user_id = "user123"
  skip_agent_check = true
}
`,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.kasm_session_status.test", "id", "user123-kasm123"),
					resource.TestCheckResourceAttr("data.kasm_session_status.test", "operational_status", "running"),
					resource.TestCheckResourceAttr("data.kasm_session_status.test", "operational_message", "Session is running"),
					resource.TestCheckResourceAttr("data.kasm_session_status.test", "kasm_url", "https://example.com/kasm/123"),
					resource.TestCheckResourceAttr("data.kasm_session_status.test", "container_ip", "10.0.0.1"),
					resource.TestCheckResourceAttr("data.kasm_session_status.test", "port", "8443"),
					resource.TestCheckResourceAttr("data.kasm_session_status.test", "image_id", "image123"),
					resource.TestCheckResourceAttr("data.kasm_session_status.test", "image_name", "Ubuntu"),
				),
			},
		},
	})
}
