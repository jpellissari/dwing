package auth_test

import (
	"jpellissari/dwing/auth"
	"os"
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

func TestLoadEmptyFile(t *testing.T) {
	c := auth.Credentials{}

	err := c.Load("non_existent_file.json")
	assert.NoError(t, err, "Loading from a non-existent file should not return an error")
	assert.Len(t, c, 0, "Credential slice should be empty after getting from a non-existent file")
}

func TestSaveCredentials(t *testing.T) {
	c := auth.Credentials{
		{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1"},
		{Environment: "env2", Username: "user2", Password: "pass2", Nickname: "nick2"},
	}

	tf, err := os.CreateTemp("", "credentials_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tf.Name())

	err = c.Save(tf.Name())
	assert.NoError(t, err, "Saving to a valid file should not return an error")
}

func TestSaveAndLoadCredentials(t *testing.T) {
	c := auth.Credentials{
		{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1"},
		{Environment: "env2", Username: "user2", Password: "pass2", Nickname: "nick2"},
	}

	tf, err := os.CreateTemp("", "credentials_*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tf.Name())

	err = c.Save(tf.Name())
	assert.NoError(t, err, "Saving to a valid file should not return an error")

	var c2 auth.Credentials
	err = c2.Load(tf.Name())
	assert.NoError(t, err, "Loading from a valid file should not return an error")
	assert.EqualValues(t, c, c2, "The credentials retrieved should match the original credentials")
}
