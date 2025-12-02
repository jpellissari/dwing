package auth_test

import (
	"jpellissari/dwing/internal/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateCredential(t *testing.T) {
	testCases := []struct {
		name       string
		credential auth.Credential
		shouldFail bool
		message    string
	}{
		{
			name: "valid_credential",
			credential: auth.Credential{
				Environment: "env1",
				Username:    "user",
				Password:    "pass",
				Nickname:    "nick",
			},
			shouldFail: false, message: "Valid credential should pass"},
		{
			name: "missing_environment",
			credential: auth.Credential{Environment: "",
				Username: "user",
				Password: "pass",
				Nickname: "nick",
			},
			shouldFail: true,
			message:    "Empty environment should be required"},

		{
			name: "missing_username",
			credential: auth.Credential{Environment: "env1",
				Username: "",
				Password: "pass",
				Nickname: "nick",
			},
			shouldFail: true,
			message:    "Empty username should be required"},

		{
			name: "missing_password",
			credential: auth.Credential{Environment: "env1",
				Username: "user",
				Password: "",
				Nickname: "nick",
			},
			shouldFail: true,
			message:    "Empty password should be required"},

		{
			name: "missing_nickname_is_valid",
			credential: auth.Credential{Environment: "env1",
				Username: "user",
				Password: "pass",
				Nickname: "",
			},
			shouldFail: false,
			message:    "Empty nickname should not be required"},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			err := tc.credential.Validate()
			if tc.shouldFail {
				assert.Error(t, err, tc.message)
			} else {
				assert.NoError(t, err, tc.message)
			}
		})
	}
}
