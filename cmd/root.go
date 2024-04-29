package cmd

import (
	"github.com/spf13/cobra"
)

func NewRootCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "integration-go",
		Short: "Run service",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(apiCmd(), cronCmd())
	return command
}
