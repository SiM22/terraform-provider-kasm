package client

import (
    "context"
    "crypto/tls"
    "net/http"
    "sync"
    "time"
)

type ClientOption func(*Client)

type Client struct {
    BaseURL     string
    HTTPClient  *http.Client
    APIKey      string
    APISecret   string
    Version     APIVersion
    rateLimiter *RateLimiter
    retryConfig *RetryConfig
    debugMode   bool
    mu          sync.RWMutex
}

type RetryConfig struct {
    MaxRetries          int
    InitialInterval     time.Duration
    MaxInterval         time.Duration
    Multiplier          float64
    RandomizationFactor float64
}

type APIVersion string

const (
    APIVersionLatest APIVersion = "latest"
    APIVersion1_0    APIVersion = "1.0"
)

// NewClient creates a new API client with options
func NewClient(baseURL, apiKey, apiSecret string, insecure bool, options ...ClientOption) *Client {
    tr := &http.Transport{
        TLSClientConfig: &tls.Config{InsecureSkipVerify: insecure},
    }

    client := &Client{
        BaseURL:     baseURL,
        APIKey:      apiKey,
        APISecret:   apiSecret,
        Version:     APIVersionLatest,
        rateLimiter: NewRateLimiter(100, time.Second),
        retryConfig: &RetryConfig{
            MaxRetries:          3,
            InitialInterval:     100 * time.Millisecond,
            MaxInterval:         10 * time.Second,
            Multiplier:          2.0,
            RandomizationFactor: 0.1,
        },
        HTTPClient: &http.Client{
            Transport: tr,
            Timeout:   time.Second * 30,
        },
    }

    // Apply options
    for _, option := range options {
        option(client)
    }

    return client
}

// Legacy support for old code
func (c *Client) doRequestLegacy(method, path string, payload map[string]interface{}) (*http.Response, error) {
    return c.doRequest(context.Background(), method, path, payload)
}
