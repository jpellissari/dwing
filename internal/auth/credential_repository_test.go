package auth_test

import (
	"encoding/json"
	"jpellissari/dwing/internal/auth"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name        string
		setupFile   func(t *testing.T, filePath string) // Setup file state
		wantCreds   auth.Credentials
		wantErr     bool
		errContains string
	}{
		{
			name: "file does not exist returns empty credentials",
			setupFile: func(t *testing.T, filePath string) {
				// Don't create file
			},
			wantCreds: auth.Credentials{},
			wantErr:   false,
		},
		{
			name: "empty file returns empty credentials",
			setupFile: func(t *testing.T, filePath string) {
				err := os.WriteFile(filePath, []byte{}, 0644)
				require.NoError(t, err)
			},
			wantCreds: auth.Credentials{},
			wantErr:   false,
		},
		{
			name: "valid JSON with single credential",
			setupFile: func(t *testing.T, filePath string) {
				creds := auth.Credentials{
					{Username: "user1", Password: "pass1", Environment: "env1"},
				}
				data, err := json.MarshalIndent(creds, "", "  ")
				require.NoError(t, err)
				err = os.WriteFile(filePath, data, 0644)
				require.NoError(t, err)
			},
			wantCreds: auth.Credentials{
				{Username: "user1", Password: "pass1", Environment: "env1"},
			},
			wantErr: false,
		},
		{
			name: "valid JSON with multiple credentials",
			setupFile: func(t *testing.T, filePath string) {
				creds := auth.Credentials{
					{Username: "user1", Password: "pass1", Environment: "env1"},
					{Username: "user2", Password: "pass2", Environment: "env1"},
					{Username: "user3", Password: "pass3", Environment: "env1"},
				}
				data, err := json.MarshalIndent(creds, "", "  ")
				require.NoError(t, err)
				err = os.WriteFile(filePath, data, 0644)
				require.NoError(t, err)
			},
			wantCreds: auth.Credentials{
				{Username: "user1", Password: "pass1", Environment: "env1"},
				{Username: "user2", Password: "pass2", Environment: "env1"},
				{Username: "user3", Password: "pass3", Environment: "env1"},
			},
			wantErr: false,
		},
		{
			name: "invalid JSON returns error",
			setupFile: func(t *testing.T, filePath string) {
				err := os.WriteFile(filePath, []byte(`{"invalid": json}`), 0644)
				require.NoError(t, err)
			},
			wantCreds:   nil,
			wantErr:     true,
			errContains: "failed to unmarshal JSON",
		},
		{
			name: "malformed JSON returns error",
			setupFile: func(t *testing.T, filePath string) {
				err := os.WriteFile(filePath, []byte(`[{"username": "test"`), 0644)
				require.NoError(t, err)
			},
			wantCreds:   nil,
			wantErr:     true,
			errContains: "failed to unmarshal JSON",
		},
		{
			name: "corrupted file returns error",
			setupFile: func(t *testing.T, filePath string) {
				err := os.WriteFile(filePath, []byte("not json at all!"), 0644)
				require.NoError(t, err)
			},
			wantCreds:   nil,
			wantErr:     true,
			errContains: "failed to unmarshal JSON",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "credentials.json")

			// Setup file state
			tt.setupFile(t, filePath)

			repo := auth.NewJSONRepository(filePath)
			creds, err := repo.GetAll()

			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, creds)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.wantCreds, creds)
			}
		})
	}
}

func TestSave(t *testing.T) {
	tests := []struct {
		name            string
		creds           auth.Credentials
		setupDir        func(t *testing.T, tmpDir string) string // Returns filePath
		wantErr         bool
		errContains     string
		verifyFile      func(t *testing.T, filePath string)
		verifyDirectory func(t *testing.T, dirPath string)
	}{
		{
			name: "new file creates file with credentials",
			creds: auth.Credentials{
				{Username: "user1", Password: "pass1", Environment: "env1"},
			},
			setupDir: func(t *testing.T, tmpDir string) string {
				return filepath.Join(tmpDir, "credentials.json")
			},
			wantErr: false,
			verifyFile: func(t *testing.T, filePath string) {
				// Verify file exists
				_, err := os.Stat(filePath)
				assert.NoError(t, err, "File should exist")

				// Verify content
				data, err := os.ReadFile(filePath)
				require.NoError(t, err)

				var loaded auth.Credentials
				err = json.Unmarshal(data, &loaded)
				require.NoError(t, err)
				assert.Len(t, loaded, 1)
			},
		},
		{
			name:  "empty credentials creates empty array",
			creds: auth.Credentials{},
			setupDir: func(t *testing.T, tmpDir string) string {
				return filepath.Join(tmpDir, "credentials.json")
			},
			wantErr: false,
			verifyFile: func(t *testing.T, filePath string) {
				data, err := os.ReadFile(filePath)
				require.NoError(t, err)
				assert.Equal(t, "[]", string(data))
			},
		},
		{
			name: "file permissions are correct",
			creds: auth.Credentials{
				{Username: "user1", Password: "pass1", Environment: "env1"},
			},
			setupDir: func(t *testing.T, tmpDir string) string {
				return filepath.Join(tmpDir, "credentials.json")
			},
			wantErr: false,
			verifyFile: func(t *testing.T, filePath string) {
				info, err := os.Stat(filePath)
				require.NoError(t, err)
				assert.Equal(t, os.FileMode(0600), info.Mode().Perm(), "File should have 0600 permissions")
			},
		},
		{
			name: "directory permissions are correct",
			creds: auth.Credentials{
				{Username: "user1", Password: "pass1", Environment: "env1"},
			},
			setupDir: func(t *testing.T, tmpDir string) string {
				return filepath.Join(tmpDir, "subdir", "credentials.json")
			},
			wantErr: false,
			verifyDirectory: func(t *testing.T, dirPath string) {
				info, err := os.Stat(dirPath)
				require.NoError(t, err)
				assert.Equal(t, os.FileMode(0755), info.Mode().Perm(), "Directory should have 0755 permissions")
			},
		},
		{
			name: "JSON format is indented",
			creds: auth.Credentials{
				{Username: "user1", Password: "pass1", Environment: "env1"},
			},
			setupDir: func(t *testing.T, tmpDir string) string {
				return filepath.Join(tmpDir, "credentials.json")
			},
			wantErr: false,
			verifyFile: func(t *testing.T, filePath string) {
				data, err := os.ReadFile(filePath)
				require.NoError(t, err)
				content := string(data)
				assert.Contains(t, content, "\n", "JSON should be formatted with newlines")
				assert.Contains(t, content, "  ", "JSON should be indented with 2 spaces")
			},
		},
		{
			name: "multiple credentials saved correctly",
			creds: auth.Credentials{
				{Username: "user1", Password: "pass1", Environment: "env1"},
				{Username: "user2", Password: "pass2", Environment: "env1"},
				{Username: "user3", Password: "pass3", Environment: "env1"},
			},
			setupDir: func(t *testing.T, tmpDir string) string {
				return filepath.Join(tmpDir, "credentials.json")
			},
			wantErr: false,
			verifyFile: func(t *testing.T, filePath string) {
				data, err := os.ReadFile(filePath)
				require.NoError(t, err)
				var loaded auth.Credentials
				err = json.Unmarshal(data, &loaded)
				require.NoError(t, err)
				assert.Len(t, loaded, 3)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := tt.setupDir(t, tmpDir)

			repo := auth.NewJSONRepository(filePath)
			err := repo.Save(tt.creds)

			if tt.wantErr {
				assert.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)

				if tt.verifyFile != nil {
					tt.verifyFile(t, filePath)
				}

				if tt.verifyDirectory != nil {
					dirPath := filepath.Dir(filePath)
					tt.verifyDirectory(t, dirPath)
				}
			}
		})
	}
}

func TestDuplicateCredentials(t *testing.T) {
	testCases := []struct {
		name       string
		cred       auth.Credential
		wantReturn bool
	}{
		{
			name: "Return true if same env and username",
			cred: auth.Credential{
				Username:    "user1",
				Environment: "env1",
				Password:    "pass1",
			},
			wantReturn: true,
		},
		{
			name: "Return false if same username and different env",
			cred: auth.Credential{
				Username:    "user1",
				Environment: "different",
				Password:    "pass1",
			},
			wantReturn: false,
		},
		{
			name: "Return false if same env and different username",
			cred: auth.Credential{
				Username:    "different",
				Environment: "env1",
				Password:    "pass1",
			},
			wantReturn: false,
		},
		{
			name: "Return false if username and env are different",
			cred: auth.Credential{
				Username:    "different",
				Environment: "different",
				Password:    "pass1",
			},
			wantReturn: false,
		},
		{
			name: "Return true if username and env are equal but nickname is different",
			cred: auth.Credential{
				Username:    "user1",
				Environment: "env1",
				Password:    "pass1",
				Nickname:    "teste",
			},
			wantReturn: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			filePath := filepath.Join(tmpDir, "credentials.json")

			repo := auth.NewJSONRepository(filePath)

			creds := auth.Credential{
				Username:    "user1",
				Environment: "env1",
				Password:    "pass1",
				Nickname:    "nick1",
			}
			repo.Add(creds)

			equal, _ := repo.CheckDuplicate(tt.cred)
			assert.Equal(t, tt.wantReturn, equal)
		})
	}
}
