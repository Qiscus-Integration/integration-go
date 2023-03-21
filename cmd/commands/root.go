package commands

import "github.com/spf13/cobra"

// NewRootCommand function returns a Cobra command that serves as the root command
// for the integration-go application.
func NewRootCommand() *cobra.Command {
	var command = &cobra.Command{
		Use:   "integration-go",
		Short: "Run service",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.HelpFunc()(cmd, args)
		},
	}

	command.AddCommand(serverCmd(), cronCmd())
	return command
}
