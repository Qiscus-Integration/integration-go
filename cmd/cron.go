package cmd

import (
	"integration-go/internal/cron"

	"github.com/spf13/cobra"
)

func cronCmd() *cobra.Command {
	var command = &cobra.Command{
		Use:   "cron",
		Short: "Run cron server",
		Run: func(cmd *cobra.Command, args []string) {
			srv := cron.NewServer()
			srv.Run()
		},
	}

	return command
}
