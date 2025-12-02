package auth_test

import (
	"jpellissari/dwing/internal/auth"
	"testing"

	"github.com/stretchr/testify/assert"
)

type FakeCredentialRepository struct {
	Credentials auth.Credentials
}

func NewFakeCredentialRepository(c auth.Credentials) *FakeCredentialRepository {
	return &FakeCredentialRepository{
		Credentials: c,
	}
}

func (r *FakeCredentialRepository) GetCredentialsByEnv(env string) (auth.Credentials, error) {
	var filteredCreds auth.Credentials
	for _, c := range r.Credentials {
		if c.Environment == env {
			filteredCreds = append(filteredCreds, c)
		}
	}

	return filteredCreds, nil
}

func (r *FakeCredentialRepository) Load() (auth.Credentials, error) {
	return r.Credentials, nil
}

func (r *FakeCredentialRepository) Save(c auth.Credentials) error {
	return nil
}

func TestListCredentials(t *testing.T) {
	testCases := []struct {
		name        string
		credentials auth.Credentials
		expectList  auth.Credentials
		envFilter   string
	}{
		{
			name:        "return empty array if empty credentials",
			credentials: auth.Credentials{},
			expectList:  auth.Credentials{},
		},
		{
			name: "return all credentials",
			credentials: auth.Credentials{
				{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1"},
				{Environment: "env1", Username: "user2", Password: "pass2", Nickname: "nick2"},
				{Environment: "env1", Username: "user3", Password: "pass3", Nickname: "nick3"},
			},
			expectList: auth.Credentials{
				{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1"},
				{Environment: "env1", Username: "user2", Password: "pass2", Nickname: "nick2"},
				{Environment: "env1", Username: "user3", Password: "pass3", Nickname: "nick3"},
			},
		},
		{
			name: "filter credentials by env if supplied",
			credentials: auth.Credentials{
				{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1"},
				{Environment: "env2", Username: "user2", Password: "pass2", Nickname: "nick2"},
				{Environment: "env2", Username: "user3", Password: "pass3", Nickname: "nick3"},
			},
			expectList: auth.Credentials{
				{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1"},
			},
			envFilter: "env1",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewFakeCredentialRepository(tc.credentials)
			service := auth.NewCredentialService(repo)

			creds, _ := service.ListCredentials(tc.envFilter)

			assert.Equal(t, tc.expectList, creds)
		})

	}

}
