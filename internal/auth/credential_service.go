package auth

import "fmt"

type CredentialService struct {
	repo CredentialRepository
}

func NewCredentialService(repo CredentialRepository) *CredentialService {
	return &CredentialService{repo: repo}
}

func (s *CredentialService) AddCredential(cred Credential) error {
	if err := cred.Validate(); err != nil {
		return fmt.Errorf("invalid credential: %w", err)
	}

	isDuplicate, err := s.repo.CheckDuplicate(cred)
	if err != nil {
		return err
	}

	if isDuplicate {
		return fmt.Errorf("credential for environment '%s' and username '%s' already exists", cred.Environment, cred.Username)
	}

	if err := s.repo.Add(cred); err != nil {
		return fmt.Errorf("failed to add credential: %w", err)
	}

	return nil
}

func (s *CredentialService) ListCredentials(env string) (Credentials, error) {
	if env != "" {
		return s.repo.GetByEnv(env)
	}

	return s.repo.GetAll()
}

func (s *CredentialService) RemoveCredential(id string) error {
	if err := s.repo.RemoveById(id); err != nil {
		return err
	}

	return nil
}
