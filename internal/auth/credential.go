package auth

import (
	"errors"
	"fmt"
)

type Credential struct {
	Environment string `json:"environment"`
	Username    string `json:"username"`
	Password    string `json:"password"`
	Nickname    string `json:"nickname"`
}

func (c *Credential) Validate() error {
	if c.Environment == "" {
		return errors.New("environment is required")
	}
	if c.Username == "" {
		return errors.New("username is required")
	}
	if c.Password == "" {
		return errors.New("password is required")
	}
	return nil
}

type Credentials []Credential

func (c *Credentials) Add(cred Credential) error {
	if err := cred.Validate(); err != nil {
		return err
	}

	*c = append(*c, cred)

	return nil
}

func (c *Credentials) Delete(i int) error {
	cl := *c
	if i <= 0 || i > len(cl) {
		return fmt.Errorf("Credential %d does not exists", i)
	}

	*c = append(cl[:i-1], cl[i:]...)

	return nil
}

func (c *Credentials) CheckDuplicate(cred Credential) bool {
	for _, existingCred := range *c {
		if existingCred.Environment == cred.Environment && existingCred.Username == cred.Username {
			return true
		}
	}
	return false
}
