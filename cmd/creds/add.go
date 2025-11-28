package creds

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCredsAddCommand() *cobra.Command {
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
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	return addCmd
}
