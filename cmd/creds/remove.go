package creds

import (
	"errors"
	"fmt"
	"jpellissari/dwing/internal/auth"
	"jpellissari/dwing/internal/config"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCredsRemoveCommand() *cobra.Command {
	var id string

	var listCmd = &cobra.Command{
		Use:     "remove <credential_id>",
		Short:   "Remove a stored credential",
		Long:    `Remove a stored credential from the dwing credential manager.`,
		Aliases: []string{"rm"},
		Example: heredoc.Doc(`
			$ dwing creds remove <credential_id>
			$ dwing creds rm <credential_id>
		`),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return fmt.Errorf("credential ID is required")
			}
			id = args[0]

			cfg, err := config.NewDefaultConfig()
			if err != nil {
				return fmt.Errorf("failed to load config: %w", err)
			}

			repo := auth.NewJSONRepository(cfg.CredentialsPath)
			service := auth.NewCredentialService(repo)

			if err := service.RemoveCredential(id); err != nil {
				if errors.Is(err, auth.ErrCredentialNotFound) {
					fmt.Printf("❌ Credential with ID '%s' not found\n", id)
					return nil
				}
				return fmt.Errorf("failed to add credential: %w", err)
			}

			fmt.Println("Credential removed successfully")

			return nil
		},
	}

	return listCmd
}
