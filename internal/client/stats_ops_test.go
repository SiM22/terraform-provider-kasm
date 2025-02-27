//go:build unit
// +build unit

package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_GetFrameStats(t *testing.T) {
	testCases := []struct {
		name           string
		kasmID         string
		userID         string
		serverResponse interface{}
		statusCode     int
		expectError    bool
	}{
		{
			name:       "successful request",
			kasmID:     "test-kasm-id",
			userID:     "test-user-id",
			statusCode: http.StatusOK,
			serverResponse: FrameStatsResponse{
				Frame: FrameStats{
					ResX:    1920,
					ResY:    1080,
					Changed: 100,
					Clients: []FrameStatsClient{
						{
							Client:     "test-client",
							ClientTime: 123,
							Ping:       5,
						},
					},
				},
			},
			expectError: false,
		},
		{
			name:           "server error",
			kasmID:         "test-kasm-id",
			userID:         "test-user-id",
			statusCode:     http.StatusInternalServerError,
			serverResponse: map[string]interface{}{"error": "internal server error"},
			expectError:    true,
		},
		{
			name:           "invalid response",
			kasmID:         "test-kasm-id",
			userID:         "test-user-id",
			statusCode:     http.StatusOK,
			serverResponse: "invalid json",
			expectError:    true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check request method and path
				assert.Equal(t, http.MethodPost, r.Method)
				assert.Equal(t, "/api/public/frame_stats", r.URL.Path)

				// Decode request body
				var reqBody FrameStatsRequest
				err := json.NewDecoder(r.Body).Decode(&reqBody)
				assert.NoError(t, err)

				// Check request parameters
				assert.Equal(t, "test-api-key", reqBody.APIKey)
				assert.Equal(t, "test-api-secret", reqBody.APISecret)
				assert.Equal(t, tc.kasmID, reqBody.KasmID)
				assert.Equal(t, tc.userID, reqBody.UserID)

				// Set response status code
				w.WriteHeader(tc.statusCode)

				// Write response body
				if tc.statusCode == http.StatusOK {
					if responseStr, ok := tc.serverResponse.(string); ok {
						w.Write([]byte(responseStr))
					} else {
						json.NewEncoder(w).Encode(tc.serverResponse)
					}
				} else {
					json.NewEncoder(w).Encode(tc.serverResponse)
				}
			}))
			defer server.Close()

			client := &Client{
				BaseURL:    server.URL,
				APIKey:     "test-api-key",
				APISecret:  "test-api-secret",
				HTTPClient: server.Client(),
			}

			response, err := client.GetFrameStats(tc.kasmID, tc.userID)
			if tc.expectError {
				assert.Error(t, err)
				assert.Nil(t, response)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, response)
				assert.Equal(t, 1920, response.Frame.ResX)
				assert.Equal(t, 1080, response.Frame.ResY)
				assert.Equal(t, 100, response.Frame.Changed)
				assert.Len(t, response.Frame.Clients, 1)
				assert.Equal(t, "test-client", response.Frame.Clients[0].Client)
				assert.Equal(t, 123, response.Frame.Clients[0].ClientTime)
				assert.Equal(t, 5, response.Frame.Clients[0].Ping)
			}
		})
	}
}
