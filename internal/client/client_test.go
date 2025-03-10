//go:build unit
// +build unit

package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupTestServer(t *testing.T, path string, expectedMethod string, expectedBody map[string]interface{}, response interface{}, statusCode int) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check request method
		if r.Method != expectedMethod {
			t.Errorf("Expected method %s, got %s", expectedMethod, r.Method)
		}

		// Check request path
		if r.URL.Path != path {
			t.Errorf("Expected path %s, got %s", path, r.URL.Path)
		}

		// Verify request body
		var requestBody map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			t.Fatalf("Failed to decode request body: %v", err)
		}

		// Check API credentials
		if requestBody["api_key"] != "test-key" || requestBody["api_key_secret"] != "test-secret" {
			t.Error("Missing or invalid API credentials in request")
		}

		// Set response
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if response != nil {
			responseBytes, err := json.Marshal(response)
			if err != nil {
				t.Fatalf("Failed to marshal response: %v", err)
			}
			if _, err := w.Write(responseBytes); err != nil {
				t.Fatal(err)
			}
		}
	}))
}

func TestLogoutUser(t *testing.T) {
	testCases := []struct {
		name        string
		userID      string
		statusCode  int
		response    interface{}
		expectError bool
	}{
		{
			name:        "successful logout",
			userID:      "test-user-id",
			statusCode:  http.StatusOK,
			response:    map[string]interface{}{"success": true},
			expectError: false,
		},
		{
			name:        "user not found",
			userID:      "invalid-user",
			statusCode:  http.StatusNotFound,
			response:    map[string]interface{}{"error": "User not found"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedBody := map[string]interface{}{
				"target_user": map[string]string{
					"user_id": tc.userID,
				},
			}

			server := setupTestServer(t, "/api/public/logout_user", http.MethodPost, expectedBody, tc.response, tc.statusCode)
			defer server.Close()

			client := NewClient(server.URL, "test-key", "test-secret", true)
			err := client.LogoutUser(tc.userID)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestGetUserAttributes(t *testing.T) {
	testCases := []struct {
		name          string
		userID        string
		statusCode    int
		response      interface{}
		expectError   bool
		expectedAttrs map[string]interface{}
	}{
		{
			name:       "successful attributes fetch",
			userID:     "test-user-id",
			statusCode: http.StatusOK,
			response: map[string]interface{}{
				"user_attributes": map[string]interface{}{
					"user_id": "test-user-id",
					"attributes": map[string]interface{}{
						"theme":    "dark",
						"language": "en",
					},
				},
			},
			expectError: false,
			expectedAttrs: map[string]interface{}{
				"theme":    "dark",
				"language": "en",
			},
		},
		{
			name:          "user not found",
			userID:        "invalid-user",
			statusCode:    http.StatusNotFound,
			response:      map[string]interface{}{"error": "User not found"},
			expectError:   true,
			expectedAttrs: nil,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedBody := map[string]interface{}{
				"target_user": map[string]string{
					"user_id": tc.userID,
				},
			}

			server := setupTestServer(t, "/api/public/get_attributes", http.MethodPost, expectedBody, tc.response, tc.statusCode)
			defer server.Close()

			client := NewClient(server.URL, "test-key", "test-secret", true)
			attrs, err := client.GetUserAttributes(tc.userID)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
			if err == nil && tc.expectedAttrs != nil {
				if attrs.UserID != tc.userID {
					t.Errorf("Expected user_id %s, got %s", tc.userID, attrs.UserID)
				}
				for k, v := range tc.expectedAttrs {
					if attrs.Attributes[k] != v {
						t.Errorf("Expected attribute %s to be %v, got %v", k, v, attrs.Attributes[k])
					}
				}
			}
		})
	}
}

func TestUpdateUserAttributes(t *testing.T) {
	testCases := []struct {
		name        string
		userID      string
		attributes  map[string]interface{}
		statusCode  int
		response    interface{}
		expectError bool
	}{
		{
			name:   "successful update",
			userID: "test-user-id",
			attributes: map[string]interface{}{
				"theme":    "light",
				"language": "fr",
			},
			statusCode:  http.StatusOK,
			response:    map[string]interface{}{"success": true},
			expectError: false,
		},
		{
			name:   "user not found",
			userID: "invalid-user",
			attributes: map[string]interface{}{
				"theme": "dark",
			},
			statusCode:  http.StatusNotFound,
			response:    map[string]interface{}{"error": "User not found"},
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			expectedBody := map[string]interface{}{
				"target_user": map[string]interface{}{
					"user_id":    tc.userID,
					"attributes": tc.attributes,
				},
			}

			server := setupTestServer(t, "/api/public/update_user_attributes", http.MethodPost, expectedBody, tc.response, tc.statusCode)
			defer server.Close()

			client := NewClient(server.URL, "test-key", "test-secret", true)
			err := client.UpdateUserAttributes(tc.userID, tc.attributes)

			if tc.expectError && err == nil {
				t.Error("Expected error but got none")
			}
			if !tc.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
			}
		})
	}
}

func TestClient_Keepalive(t *testing.T) {
	tests := []struct {
		name         string
		kasmID       string
		mockResponse string
		statusCode   int
		expectError  bool
	}{
		{
			name:         "Successful keepalive",
			kasmID:       "test-kasm-id",
			mockResponse: `{"usage_reached": false}`,
			statusCode:   http.StatusOK,
			expectError:  false,
		},
		{
			name:         "Keepalive with usage reached",
			kasmID:       "test-kasm-id",
			mockResponse: `{"usage_reached": true}`,
			statusCode:   http.StatusOK,
			expectError:  false,
		},
		{
			name:         "Keepalive with invalid response",
			kasmID:       "test-kasm-id",
			mockResponse: `invalid`,
			statusCode:   http.StatusOK,
			expectError:  true,
		},
		{
			name:         "Keepalive with server error",
			kasmID:       "test-kasm-id",
			mockResponse: `{}`,
			statusCode:   http.StatusInternalServerError,
			expectError:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(tt.statusCode)
				w.Write([]byte(tt.mockResponse))
			}))
			defer server.Close()

			client := &Client{
				HTTPClient: server.Client(),
				BaseURL:    server.URL,
				APIKey:     "test-api-key",
				APISecret:  "test-api-secret",
			}

			response, err := client.Keepalive(tt.kasmID)
			if (err != nil) != tt.expectError {
				t.Errorf("Client.Keepalive() error = %v, expectError %v", err, tt.expectError)
				return
			}

			if !tt.expectError && response == nil {
				t.Error("Client.Keepalive() returned nil response")
			}
		})
	}
}
