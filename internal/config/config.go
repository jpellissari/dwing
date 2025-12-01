package config

import (
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	CredentialsPath string `json:"credentials_path"`
}

func NewDefaultConfig() (*Config, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	return NewConfig(filepath.Join(homeDir, ".dwing", "credentials.json"))

}

func NewConfig(credentialsPath string) (*Config, error) {
	cfg := &Config{
		CredentialsPath: credentialsPath,
	}

	err := cfg.Validate()
	if err != nil {
		return nil, err
	}

	err = cfg.EnsureCredentialsDirExists()
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) Validate() error {
	if c.CredentialsPath == "" {
		return fmt.Errorf("credentials path cannot be empty")
	}

	if !filepath.IsAbs(c.CredentialsPath) {
		return fmt.Errorf("credentials path must be an absolute path: %s", c.CredentialsPath)
	}

	return nil
}

func (c *Config) EnsureCredentialsDirExists() error {
	dir := filepath.Dir(c.CredentialsPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create credentials directory: %w", err)
	}
	return nil
}
