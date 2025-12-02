package creds

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCredsCmd() *cobra.Command {
	var credsCmd = &cobra.Command{
		Use:   "creds <command> [flags]",
		Short: "Manage your credentials",
		Long:  `Manage your credentials and generate tokens as needed for various environments.`,
		Example: heredoc.Doc(`
			$ dwing creds ls
			$ dwing creds add
			$ dwing creds rm <credential-id>
			$ dwing creds login <credential-id>
		`),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	credsGroup := cobra.Group{
		ID:    "creds",
		Title: "Credentials Management",
	}
	credsCmd.AddGroup(&credsGroup)

	credsAddCmd := NewCredsAddCommand()
	credsAddCmd.GroupID = credsGroup.ID

	credsListCmd := NewCredsListCommand()
	credsListCmd.GroupID = credsGroup.ID

	credsCmd.AddCommand(credsAddCmd)
	credsCmd.AddCommand(credsListCmd)

	return credsCmd
}
