package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/LeartS/mb/internal/config"
)

// Client is an authenticated HTTP client for the Metabase API.
type Client struct {
	baseURL      string
	apiKey       string
	sessionToken string
	http         *http.Client
}

// New creates a Client from the current configuration.
func New() (*Client, error) {
	cfg, err := config.Load()
	if err != nil {
		return nil, err
	}
	host, err := cfg.ActiveHost()
	if err != nil {
		return nil, err
	}
	apiKey, sessionToken, err := cfg.ActiveAuth()
	if err != nil {
		return nil, err
	}
	return &Client{
		baseURL:      strings.TrimRight(host, "/"),
		apiKey:       apiKey,
		sessionToken: sessionToken,
		http: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

// NewWithCredentials creates a Client with explicit credentials (used during auth login).
func NewWithCredentials(host, apiKey, sessionToken string) *Client {
	return &Client{
		baseURL:      strings.TrimRight(host, "/"),
		apiKey:       apiKey,
		sessionToken: sessionToken,
		http: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) url(path string) string {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	return c.baseURL + path
}

func (c *Client) setAuth(req *http.Request) {
	if c.apiKey != "" {
		req.Header.Set("X-Api-Key", c.apiKey)
	} else if c.sessionToken != "" {
		req.Header.Set("X-Metabase-Session", c.sessionToken)
	}
}

// Do executes a raw HTTP request and returns the response body as bytes.
func (c *Client) Do(method, path string, body io.Reader) ([]byte, int, error) {
	req, err := http.NewRequest(method, c.url(path), body)
	if err != nil {
		return nil, 0, fmt.Errorf("creating request: %w", err)
	}
	c.setAuth(req)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("executing request: %w", err)
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, fmt.Errorf("reading response: %w", err)
	}
	return data, resp.StatusCode, nil
}

// Get performs a GET request and returns the raw JSON bytes.
func (c *Client) Get(path string) (json.RawMessage, error) {
	data, status, err := c.Do("GET", path, nil)
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("GET %s: HTTP %d: %s", path, status, string(data))
	}
	return json.RawMessage(data), nil
}

// Post performs a POST request with a JSON body.
func (c *Client) Post(path string, payload any) (json.RawMessage, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling request body: %w", err)
	}
	data, status, err := c.Do("POST", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("POST %s: HTTP %d: %s", path, status, string(data))
	}
	return json.RawMessage(data), nil
}

// Put performs a PUT request with a JSON body.
func (c *Client) Put(path string, payload any) (json.RawMessage, error) {
	body, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("marshaling request body: %w", err)
	}
	data, status, err := c.Do("PUT", path, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("PUT %s: HTTP %d: %s", path, status, string(data))
	}
	return json.RawMessage(data), nil
}

// Delete performs a DELETE request.
func (c *Client) Delete(path string) (json.RawMessage, error) {
	data, status, err := c.Do("DELETE", path, nil)
	if err != nil {
		return nil, err
	}
	if status < 200 || status >= 300 {
		return nil, fmt.Errorf("DELETE %s: HTTP %d: %s", path, status, string(data))
	}
	return json.RawMessage(data), nil
}

// DoRaw performs a request with a raw body (for the api escape hatch).
func (c *Client) DoRaw(method, path string, body []byte) ([]byte, int, error) {
	var reader io.Reader
	if len(body) > 0 {
		reader = bytes.NewReader(body)
	}
	return c.Do(method, path, reader)
}
