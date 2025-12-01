package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewConfig(t *testing.T) {
	t.Run("Valid absolute path", func(t *testing.T) {
		tmpDir := t.TempDir()
		credPath := filepath.Join(tmpDir, "credentials.json")

		cfg, err := NewConfig(credPath)

		require.NoError(t, err)
		assert.Equal(t, credPath, cfg.CredentialsPath)

		// Verify directory was created
		dir := filepath.Dir(credPath)
		_, err = os.Stat(dir)
		assert.NoError(t, err, "Directory should exist")
	})

	t.Run("Empty path should fail validation", func(t *testing.T) {
		cfg, err := NewConfig("")

		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "cannot be empty")
	})

	t.Run("Relative path should fail validation", func(t *testing.T) {
		cfg, err := NewConfig("./credentials.json")

		assert.Error(t, err)
		assert.Nil(t, cfg)
		assert.Contains(t, err.Error(), "must be an absolute path")
	})
}

func TestValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     *Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "Valid config",
			cfg: &Config{
				CredentialsPath: "/home/user/.dwing/credentials.json",
			},
			wantErr: false,
		},
		{
			name: "Empty path",
			cfg: &Config{
				CredentialsPath: "",
			},
			wantErr: true,
			errMsg:  "cannot be empty",
		},
		{
			name: "Relative path",
			cfg: &Config{
				CredentialsPath: "relative/path.json",
			},
			wantErr: true,
			errMsg:  "must be an absolute path",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestEnsureCredentialsDirExists(t *testing.T) {
	t.Run("Creates directory successfully", func(t *testing.T) {
		tmpDir := t.TempDir()
		credPath := filepath.Join(tmpDir, "subdir", "credentials.json")

		cfg := &Config{CredentialsPath: credPath}
		err := cfg.EnsureCredentialsDirExists()

		require.NoError(t, err)

		// Verify directory exists
		dir := filepath.Dir(credPath)
		info, err := os.Stat(dir)
		require.NoError(t, err)
		assert.True(t, info.IsDir())
	})
}
