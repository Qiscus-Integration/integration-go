package commands

import (
	"integration-go/cron"

	"github.com/spf13/cobra"
)

// The cronCmd returns a new instance of cobra.Command that represents the "cron" command
// for run the cron job
func cronCmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "cron",
		Short: "Run cron",
		Run: func(cmd *cobra.Command, args []string) {
			cron := cron.NewCron()
			cron.Run()
		},
	}

	return command
}
