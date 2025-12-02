package auth

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

type CredentialRepository interface {
	Load() (Credentials, error)
	Save(c Credentials) error
	GetById(id string) (Credential, error)
	GetByEnv(env string) (Credentials, error)
	RemoveById(id string) error
}

type JSONRepository struct {
	filePath string
}

func NewJSONRepository(filePath string) *JSONRepository {
	return &JSONRepository{filePath: filePath}
}

func (r *JSONRepository) GetById(id string) (Credential, error) {
	creds, err := r.Load()
	if err != nil {
		return Credential{}, err
	}

	for _, c := range creds {
		if c.ID == id {
			return c, nil
		}
	}

	return Credential{}, fmt.Errorf("credential with ID '%s' not found", id)
}

func (r *JSONRepository) GetByEnv(env string) (Credentials, error) {
	if env == "" {
		return nil, fmt.Errorf("environment not informed!")
	}

	creds, err := r.Load()
	if err != nil {
		return nil, err
	}

	var filteredCreds Credentials
	for _, c := range creds {
		if c.Environment == env {
			filteredCreds = append(filteredCreds, c)
		}
	}

	return filteredCreds, nil
}

func (r *JSONRepository) RemoveById(id string) error {
	creds, err := r.Load()
	if err != nil {
		return err
	}

	var updatedCreds Credentials
	var found bool
	for i, c := range creds {
		if c.ID == id {
			found = true
			updatedCreds = append(creds[:i], creds[i+1:]...)
		}
	}

	if !found {
		return ErrCredentialNotFound
	}

	if err := r.Save(updatedCreds); err != nil {
		return fmt.Errorf("failed to save updated credentials: %w", err)
	}

	return nil
}

func (r *JSONRepository) Load() (Credentials, error) {
	if _, err := os.Stat(r.filePath); errors.Is(err, os.ErrNotExist) {
		return Credentials{}, nil
	}

	data, err := os.ReadFile(r.filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	if len(data) == 0 {
		return Credentials{}, nil
	}

	var c Credentials
	if err := json.Unmarshal(data, &c); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return c, nil
}

func (r *JSONRepository) Save(c Credentials) error {
	dir := filepath.Dir(r.filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal JSON: %w", err)
	}

	if err := os.WriteFile(r.filePath, data, 0600); err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}

	return nil
}
