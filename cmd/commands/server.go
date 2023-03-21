package commands

import (
	"integration-go/server"

	"github.com/spf13/cobra"
)

// The serverCmd function returns a Cobra command that starts a server on the specified port
// or the default port if no port is specified, defaults to 8080
func serverCmd() *cobra.Command {
	var port int
	var command = &cobra.Command{
		Use:   "server",
		Short: "Run server",
		Run: func(cmd *cobra.Command, args []string) {
			srv := server.NewServer()
			srv.Run(port)
		},
	}

	command.Flags().IntVar(&port, "port", 8080, "Listen on given port")
	return command
}
