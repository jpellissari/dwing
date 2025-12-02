package creds

import (
	"fmt"
	"jpellissari/dwing/internal/auth"
	"jpellissari/dwing/internal/config"
	"os"

	"github.com/MakeNowJust/heredoc"
	"github.com/olekukonko/tablewriter"
	"github.com/spf13/cobra"
)

func NewCredsListCommand() *cobra.Command {
	var env string

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

			creds, err := service.ListCredentials(env)
			if err != nil {
				return fmt.Errorf("failed to list credentials: %w", err)
			}

			renderTable(creds)

			return nil
		},
	}

	listCmd.Flags().StringVarP(&env, "env", "e", "", "Filter credentials by environment")

	return listCmd
}

func renderTable(creds auth.Credentials) {
	if len(creds) == 0 {
		fmt.Println("No credentials found.")
		fmt.Println("Try adding some with 'dwing creds add'")
		return
	}

	header := []string{"ID", "Environment", "Username", "Nickname"}

	data := [][]string{}
	for _, c := range creds {
		row := []string{c.ID, c.Environment, c.Username, c.Nickname}
		data = append(data, row)
	}

	table := tablewriter.NewTable(os.Stdout)
	table.Header(header)
	table.Bulk(data)
	table.Render()
}
