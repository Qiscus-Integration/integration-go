package cmd

import (
	"integration-go/cron"

	"github.com/spf13/cobra"
)

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
