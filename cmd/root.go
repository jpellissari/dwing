package cmd

import (
	"jpellissari/dwing/cmd/creds"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
)

func NewCmdRoot() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:   "dwing <command> <subcommand> [flags]",
		Short: "Your Developer Wingman CLI",
		Long:  `Dwing is your developer wingman, designed to make you faster on your day-to-day tasks.`,
		Example: heredoc.Doc(`
			$dwing creds ls
		`),
		Run: func(cmd *cobra.Command, args []string) {
			_ = cmd.Help()
		},
	}

	rootCmd.AddCommand(creds.NewCredsCmd())

	return rootCmd
}
