package creds

import (
	"fmt"
	"jpellissari/dwing/internal/auth"
	"jpellissari/dwing/internal/config"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCredsListCommand() *cobra.Command {
	var listCmd = &cobra.Command{
		Use:     "list",
		Short:   "List all stored credentials",
		Long:    `List all stored credentials in the dwing credential manager.`,
		Aliases: []string{"ls"},
		Example: heredoc.Doc(`
			$ dwing creds list [--env <environment>]
			$ dwing creds ls [-e <environment>]
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := config.NewDefaultConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			repo := auth.NewJSONRepository(cfg.CredentialsPath)
			service := auth.NewCredentialService(repo)

			creds, err := service.ListCredentials("teste")
			if err != nil {
				return fmt.Errorf("failed to add credential: %w", err)
			}

			fmt.Printf("creds %+v\n", creds)
			return nil
		},
	}

	return listCmd
}
