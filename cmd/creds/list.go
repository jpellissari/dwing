package creds

import (
	"errors"
	"fmt"
	"jpellissari/dwing/internal/auth"
	"jpellissari/dwing/internal/config"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCredsListCommand() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:   "list [flags]",
		Short: "List all stored credentials",
		Long:  `List all stored credentials in the dwing credential manager.`,
		Example: heredoc.Doc(`
			$ dwing creds list|ls
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}

	return listCmd
}

func listCredential(c *auth.Credential) error {
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
