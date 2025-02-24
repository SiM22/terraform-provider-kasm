package client

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetLicenses(t *testing.T) {
	testCases := []struct {
		name           string
		responseStatus int
		responseBody   interface{}
		expectError    bool
		expectedError  string
		expectedCount  int
	}{
		{
			name:           "success",
			responseStatus: http.StatusOK,
			responseBody: GetLicensesResponse{
				Success: true,
				Code:    200,
				Message: "Success",
				Data: []License{
					{
						LicenseID:   "test-license-1",
						Expiration:  "2025-01-01",
						IssuedAt:    "2024-01-01",
						Limit:       10,
						IsVerified:  true,
						LicenseType: "enterprise",
						SKU:         "ENT-10",
						Features: LicenseFeatures{
							AutoScaling:    true,
							Branding:       true,
							SessionStaging: true,
							SessionCasting: true,
						},
					},
					{
						LicenseID:   "test-license-2",
						Expiration:  "2025-01-01",
						IssuedAt:    "2024-01-01",
						Limit:       5,
						IsVerified:  true,
						LicenseType: "standard",
						SKU:         "STD-5",
						Features: LicenseFeatures{
							AutoScaling:    false,
							Branding:       false,
							SessionStaging: true,
							SessionCasting: true,
						},
					},
				},
			},
			expectError:   false,
			expectedCount: 2,
		},
		{
			name:           "api error",
			responseStatus: http.StatusOK,
			responseBody: GetLicensesResponse{
				Success: false,
				Code:    400,
				Message: "Invalid API key",
			},
			expectError:   true,
			expectedError: "API request was not successful: Invalid API key",
			expectedCount: 0,
		},
		{
			name:           "http error",
			responseStatus: http.StatusInternalServerError,
			responseBody:   "Internal Server Error",
			expectError:    true,
			expectedError:  "API request failed with status code: 500",
			expectedCount:  0,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Verify request
				assert.Equal(t, "/api/public/get_licenses", r.URL.Path)
				assert.Equal(t, "POST", r.Method)
				assert.Equal(t, "application/json", r.Header.Get("Content-Type"))

				// Send response
				w.WriteHeader(tc.responseStatus)
				if tc.responseStatus == http.StatusOK {
					json.NewEncoder(w).Encode(tc.responseBody)
				} else {
					w.Write([]byte(tc.responseBody.(string)))
				}
			}))
			defer server.Close()

			client := &Client{
				HTTPClient: server.Client(),
				BaseURL:    server.URL,
				APIKey:     "test-key",
				APISecret:  "test-secret",
			}

			licenses, err := client.GetLicenses()

			if tc.expectError {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tc.expectedError)
				assert.Nil(t, licenses)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, licenses)
				assert.Equal(t, tc.expectedCount, len(licenses))

				if tc.expectedCount > 0 {
					// Verify first license details
					assert.Equal(t, "test-license-1", licenses[0].LicenseID)
					assert.Equal(t, "enterprise", licenses[0].LicenseType)
					assert.True(t, licenses[0].Features.AutoScaling)
				}
			}
		})
	}
}
