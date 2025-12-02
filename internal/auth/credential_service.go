package auth

import "fmt"

type CredentialService struct {
	repo CredentialRepository
}

func NewCredentialService(repo CredentialRepository) *CredentialService {
	return &CredentialService{repo: repo}
}

func (s *CredentialService) AddCredential(cred *Credential) error {
	if err := cred.Validate(); err != nil {
		return fmt.Errorf("invalid credential: %w", err)
	}

	creds, err := s.repo.Load()
	if err != nil {
		return fmt.Errorf("failed to load credentials: %w", err)
	}

	if creds.CheckDuplicate(*cred) {
		return fmt.Errorf("credential for environment '%s' and username '%s' already exists", cred.Environment, cred.Username)
	}

	if err := creds.Add(*cred); err != nil {
		return fmt.Errorf("failed to add credential: %w", err)
	}

	if err := s.repo.Save(creds); err != nil {
		return fmt.Errorf("failed to save credentials: %w", err)
	}

	return nil
}

func (s *CredentialService) ListCredentials(env string) (Credentials, error) {
	if env != "" {
		return s.repo.GetCredentialsByEnv(env)
	}

	return s.repo.Load()
}
