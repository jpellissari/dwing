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

func (r *FakeCredentialRepository) Add(cred auth.Credential) error {
	r.Credentials = append(r.Credentials, cred)
	return nil
}

func (r *FakeCredentialRepository) CheckDuplicate(cred auth.Credential) (bool, error) {
	for _, c := range r.Credentials {
		if c.Environment == cred.Environment && c.Username == cred.Username && c.Nickname == cred.Nickname {
			return true, nil
		}
	}
	return false, nil
}

func (r *FakeCredentialRepository) GetById(id string) (auth.Credential, error) {
	return auth.Credential{}, nil
}

func (r *FakeCredentialRepository) GetByEnv(env string) (auth.Credentials, error) {
	var filteredCreds auth.Credentials
	for _, c := range r.Credentials {
		if c.Environment == env {
			filteredCreds = append(filteredCreds, c)
		}
	}

	return filteredCreds, nil
}

func (r *FakeCredentialRepository) RemoveById(id string) error {
	for i, c := range r.Credentials {
		if c.ID == id {
			r.Credentials = append(r.Credentials[:i], r.Credentials[i+1:]...)
			return nil
		}
	}
	return auth.ErrCredentialNotFound
}

func (r *FakeCredentialRepository) GetAll() (auth.Credentials, error) {
	return r.Credentials, nil
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

func TestRemoveCredential(t *testing.T) {
	testCases := []struct {
		name        string
		credentials auth.Credentials
		credId      string
		expectError bool
	}{
		{
			name: "remove existing credential",
			credentials: auth.Credentials{
				{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1", ID: "1"},
				{Environment: "env2", Username: "user2", Password: "pass2", Nickname: "nick2", ID: "2"},
			},
			credId:      "1",
			expectError: false,
		},
		{
			name: "error when removing non-existing credential",
			credentials: auth.Credentials{
				{Environment: "env1", Username: "user1", Password: "pass1", Nickname: "nick1", ID: "1"},
			},
			credId:      "2",
			expectError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			repo := NewFakeCredentialRepository(tc.credentials)
			service := auth.NewCredentialService(repo)

			err := service.RemoveCredential(tc.credId)

			if tc.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				remainingCreds, _ := repo.GetAll()
				for _, cred := range remainingCreds {
					assert.NotEqual(t, tc.credId, cred.ID)
				}
			}
		})
	}
}
