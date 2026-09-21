// Package typesafe is an unofficial Go client for the TypeSafe AI API.
package typesafe

import (
	"errors"
	"net/http"
	"os"
	"time"
)

const (
	defaultBaseURL    = "https://api.typesafe.ai"
	defaultModel      = "jev-latest"
	defaultMaxRetries = 3
)

// Client talks to the TypeSafe AI API.
type Client struct {
	apiKey     string
	baseURL    string
	model      string
	httpClient *http.Client
	maxRetries int
	cache      *ttlCache
}

// Option configures a Client.
type Option func(*Client)

// WithAPIKey sets the API key, overriding TYPESAFE_API_KEY.
func WithAPIKey(key string) Option { return func(c *Client) { c.apiKey = key } }

// WithBaseURL sets the base URL, overriding TYPESAFE_BASE_URL.
func WithBaseURL(url string) Option { return func(c *Client) { c.baseURL = url } }

// WithModel sets the default model, overriding TYPESAFE_DEFAULT_MODEL.
func WithModel(model string) Option { return func(c *Client) { c.model = model } }

// WithHTTPClient sets the underlying *http.Client.
func WithHTTPClient(hc *http.Client) Option { return func(c *Client) { c.httpClient = hc } }

// WithMaxRetries sets how many times a retryable failure (network error,
// 429, or 5xx) is retried, with backoff honoring any Retry-After header.
// Default 3.
func WithMaxRetries(n int) Option { return func(c *Client) { c.maxRetries = n } }

// WithCache enables an opt-in, in-memory cache for SystemOne responses,
// keyed by the request body and evicted after ttl. Disabled by default.
func WithCache(maxEntries int, ttl time.Duration) Option {
	return func(c *Client) { c.cache = newTTLCache(maxEntries, ttl) }
}

// NewClient builds a Client from environment variables and options.
// It fails if no API key is configured.
func NewClient(opts ...Option) (*Client, error) {
	c := &Client{
		apiKey:     os.Getenv("TYPESAFE_API_KEY"),
		baseURL:    defaultBaseURL,
		model:      defaultModel,
		httpClient: http.DefaultClient,
		maxRetries: defaultMaxRetries,
	}
	if v := os.Getenv("TYPESAFE_BASE_URL"); v != "" {
		c.baseURL = v
	}
	if v := os.Getenv("TYPESAFE_DEFAULT_MODEL"); v != "" {
		c.model = v
	}
	for _, opt := range opts {
		opt(c)
	}
	if c.apiKey == "" {
		return nil, errors.New("typesafe: missing API key (set TYPESAFE_API_KEY or use WithAPIKey)")
	}
	return c, nil
}
