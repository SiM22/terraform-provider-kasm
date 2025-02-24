// New file: internal/client/http.go
package client

import (
    "context"
    "encoding/json"
    "fmt"
    "net/http"
	"io"
	"bytes"
	"time"
)

func (c *Client) doRequest(ctx context.Context, method, endpoint string, payload interface{}) (*http.Response, error) {
    // Check rate limit
    if err := c.rateLimiter.Take(); err != nil {
        return nil, fmt.Errorf("rate limit: %w", err)
    }

    var body io.Reader
    if payload != nil {
        jsonData, err := json.Marshal(payload)
        if err != nil {
            return nil, fmt.Errorf("marshaling request payload: %w", err)
        }
        body = bytes.NewBuffer(jsonData)
    }

    req, err := http.NewRequestWithContext(ctx, method, c.BaseURL+endpoint, body)
    if err != nil {
        return nil, fmt.Errorf("creating request: %w", err)
    }

    // Add headers
    req.Header.Set("Content-Type", "application/json")
    if c.APIKey != "" {
        req.Header.Set("X-Api-Key", c.APIKey)
        req.Header.Set("X-Api-Secret", c.APISecret)
    }

    backoff := NewExponentialBackoff(c.retryConfig)
    var resp *http.Response

    // Retry loop
    for attempt := 0; attempt <= c.retryConfig.MaxRetries; attempt++ {
        if attempt > 0 {
            select {
            case <-ctx.Done():
                return nil, ctx.Err()
            case <-time.After(backoff.NextBackOff()):
            }
        }

        resp, err = c.HTTPClient.Do(req)
        if err != nil {
            if attempt == c.retryConfig.MaxRetries {
                return nil, fmt.Errorf("request failed after %d attempts: %w",
                    c.retryConfig.MaxRetries, err)
            }
            continue
        }

        if shouldRetry(resp.StatusCode) && attempt < c.retryConfig.MaxRetries {
            resp.Body.Close()
            continue
        }

        break
    }

    return resp, nil
}

func shouldRetry(statusCode int) bool {
    return statusCode == 429 || // Rate limit exceeded
           statusCode == 408 || // Request timeout
           statusCode == 500 || // Internal server error
           statusCode == 502 || // Bad gateway
           statusCode == 503 || // Service unavailable
           statusCode == 504    // Gateway timeout
}
