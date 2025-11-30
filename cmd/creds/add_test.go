package creds

import (
	"testing"

	"jpellissari/dwing/internal/auth"

	"github.com/stretchr/testify/assert"
)

func TestValidateRequiredFlags(t *testing.T) {
	tests := []struct {
		name    string
		cred    auth.Credential
		wantErr bool
	}{
		{
			name: "All required flags provided",
			cred: auth.Credential{
				Username:    "user1",
				Password:    "pass1",
				Environment: "dev",
			},
			wantErr: false,
		},
		{
			name: "Missing username",
			cred: auth.Credential{
				Password:    "pass1",
				Environment: "dev",
			},
			wantErr: true,
		},
		{
			name: "Missing password",
			cred: auth.Credential{
				Username:    "user1",
				Environment: "dev",
			},
			wantErr: true,
		},
		{
			name: "Missing environment",
			cred: auth.Credential{
				Username: "user1",
				Password: "pass1",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateFlags(&tt.cred)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestRequiredFieldValidator(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "Non-empty input",
			input:   "valid",
			wantErr: false,
		},
		{
			name:    "Empty input",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := requiredFieldValidator(tt.input)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}
