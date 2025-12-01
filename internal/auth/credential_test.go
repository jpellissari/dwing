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

func TestAdd(t *testing.T) {
	c := auth.Credentials{}

	expected := auth.Credential{
		Environment: "env1",
		Username:    "user",
		Password:    "pass",
		Nickname:    "nick",
	}

	err := c.Add(expected)
	assert.NoError(t, err, "Adding a valid credential should not return an error")

	assert.Len(t, c, 1, "Credential slice length should be 1 after adding a credential")
	assert.EqualValues(t, expected, c[0], "The added credential should match the expected values")
}

func TestDelete(t *testing.T) {
	c := auth.Credentials{
		{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1"},
		{Environment: "env2", Username: "user2", Password: "pass2", Nickname: "nick2"},
		{Environment: "env3", Username: "user3", Password: "pass3", Nickname: "nick3"},
	}

	c.Delete(2)

	assert.Len(t, c, 2, "Credential slice length should be 2 after deleting a credential")
	assert.EqualValues(t, "env1", c[0].Environment, "First credential should remain unchanged")
	assert.EqualValues(t, "env3", c[1].Environment, "Third credential should now be second after deletion")
}

func TestDeleteInvalidIndex(t *testing.T) {
	c := auth.Credentials{
		{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1"},
	}

	err := c.Delete(0)
	assert.Error(t, err, "Deleting with index 0 should return an error")

	err = c.Delete(2)
	assert.Error(t, err, "Deleting with an out-of-bounds index should return an error")
}

func TestDuplicateCredentials(t *testing.T) {
	c := auth.Credentials{
		{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1"},
	}

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
			equal := c.CheckDuplicate(tt.cred)
			assert.Equal(t, tt.wantReturn, equal)
		})
	}
}
