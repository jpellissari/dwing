package creds

import (
	"errors"
	"fmt"
	"jpellissari/dwing/internal/auth"
	"jpellissari/dwing/internal/config"

	"github.com/MakeNowJust/heredoc"
	"github.com/charmbracelet/huh"
	"github.com/spf13/cobra"
)

func NewCredsAddCommand() *cobra.Command {
	var cred = auth.Credential{}

	var addCmd = &cobra.Command{
		Use:   "add [flags]",
		Short: "Add a new credential",
		Long:  `Add a new credential to your credential store, either interactively or by specifying details via flags.`,
		Example: heredoc.Doc(`
			$ dwing creds add (interactive)
			$ dwing creds add -u myuser -p mypass -e dev -n mynick
		`),
		Annotations: map[string]string{
			"help:arguments": heredoc.Doc(`
				A credential can be added interactively using the 'dwing creds add' command.
				Or can be specified directly using flags:
				-u, --username <username>        Specify the username for the credential
				-p, --password <password>        Specify the password for the credential
				-e, --env <environment>          Specify the environment (e.g., dev, staging, prod)
				-n, --nickname <nickname>        Specify a nickname for easy reference (optional)
			`),
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			flagMode := cred.Username != "" || cred.Password != "" || cred.Environment != "" || cred.Nickname != ""
			if !flagMode {
				err := promptForCredential(&cred)
				if err != nil {
					return fmt.Errorf("failed to get credential input: %w", err)
				}

				return addCredential(cred)
			}

			if err := validateFlags(&cred); err != nil {
				return err
			}

			return addCredential(cred)
		},
	}

	addCmd.Flags().StringVarP(&cred.Environment, "environment", "e", "", "Environment (required)")
	addCmd.Flags().StringVarP(&cred.Username, "username", "u", "", "Username (required)")
	addCmd.Flags().StringVarP(&cred.Password, "password", "p", "", "Password (required)")
	addCmd.Flags().StringVarP(&cred.Nickname, "nickname", "n", "", "Nickname (optional)")

	return addCmd
}

func validateFlags(c *auth.Credential) error {
	allRequiredFlagsSet := c.Username != "" && c.Password != "" && c.Environment != ""

	if !allRequiredFlagsSet {
		return fmt.Errorf("when using flags, --username, --password, and --environment are required")
	}

	return nil
}

func promptForCredential(c *auth.Credential) error {
	form := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Environment").
				Prompt(">").
				Value(&c.Environment).
				Validate(requiredFieldValidator),
			huh.NewInput().
				Title("Username").
				Prompt(">").
				Value(&c.Username).
				Validate(requiredFieldValidator),
			huh.NewInput().
				Title("Password").
				Prompt(">").
				Value(&c.Password).
				EchoMode(huh.EchoModePassword).
				Validate(requiredFieldValidator),
			huh.NewInput().
				Title("Nickname").
				Prompt(">").
				Value(&c.Nickname),
		),
	)

	return form.Run()
}

func requiredFieldValidator(s string) error {
	if s == "" {
		return errors.New("this field is required")
	}
	return nil
}

func addCredential(c auth.Credential) error {
	cfg, err := config.NewDefaultConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	repo := auth.NewJSONRepository(cfg.CredentialsPath)
	service := auth.NewCredentialService(repo)

	err = service.AddCredential(c)
	if err != nil {
		return fmt.Errorf("failed to add credential: %w", err)
	}

	fmt.Printf("Credential added successfully: (%s) - %s\n", c.Environment, c.Username)

	return nil
}
