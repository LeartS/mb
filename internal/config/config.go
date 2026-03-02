package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the CLI configuration.
type Config struct {
	DefaultHost string          `json:"default_host,omitempty"`
	Hosts       map[string]Host `json:"hosts,omitempty"`
}

// Host holds per-instance configuration.
type Host struct {
	APIKey       string `json:"api_key,omitempty"`
	SessionToken string `json:"session_token,omitempty"`
}

// Dir returns the configuration directory path.
func Dir() (string, error) {
	if d := os.Getenv("MB_CONFIG_DIR"); d != "" {
		return d, nil
	}
	home, err := os.UserConfigDir()
	if err != nil {
		return "", fmt.Errorf("determining config directory: %w", err)
	}
	return filepath.Join(home, "mb"), nil
}

// Path returns the full path to the config file.
func Path() (string, error) {
	dir, err := Dir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "config.json"), nil
}

// Load reads the config from disk. Returns a zero Config if the file doesn't exist.
func Load() (*Config, error) {
	p, err := Path()
	if err != nil {
		return nil, err
	}
	data, err := os.ReadFile(p)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return &Config{Hosts: make(map[string]Host)}, nil
		}
		return nil, fmt.Errorf("reading config: %w", err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parsing config: %w", err)
	}
	if cfg.Hosts == nil {
		cfg.Hosts = make(map[string]Host)
	}
	return &cfg, nil
}

// Save writes the config to disk, creating the directory if needed.
func (c *Config) Save() error {
	p, err := Path()
	if err != nil {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(p), 0o700); err != nil {
		return fmt.Errorf("creating config directory: %w", err)
	}
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	return os.WriteFile(p, data, 0o600)
}

// ActiveHost returns the host URL to use, respecting environment overrides.
// Priority: MB_HOST env > DefaultHost in config.
func (c *Config) ActiveHost() (string, error) {
	if h := os.Getenv("MB_HOST"); h != "" {
		return h, nil
	}
	if c.DefaultHost != "" {
		return c.DefaultHost, nil
	}
	return "", fmt.Errorf("no Metabase host configured; run 'mb auth login' or set MB_HOST")
}

// ActiveAuth returns the API key or session token for the active host.
// Priority: MB_API_KEY env > MB_SESSION_TOKEN env > host config.
func (c *Config) ActiveAuth() (apiKey, sessionToken string, err error) {
	if k := os.Getenv("MB_API_KEY"); k != "" {
		return k, "", nil
	}
	if t := os.Getenv("MB_SESSION_TOKEN"); t != "" {
		return "", t, nil
	}
	host, err := c.ActiveHost()
	if err != nil {
		return "", "", err
	}
	h, ok := c.Hosts[host]
	if !ok {
		return "", "", fmt.Errorf("no credentials for host %s; run 'mb auth login'", host)
	}
	return h.APIKey, h.SessionToken, nil
}
