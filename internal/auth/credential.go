package auth

import (
	"errors"
)

type Credential struct {
	ID          string `json:"id"`
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
