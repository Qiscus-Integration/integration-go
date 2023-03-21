package main

import (
	"integration-go/cmd/commands"

	_ "github.com/joho/godotenv/autoload"

	"github.com/rs/zerolog/log"
)

// The main function initializes the root command and executes it.
func main() {
	if err := commands.NewRootCommand().Execute(); err != nil {
		log.Fatal().Msgf("failed run app: %s", err.Error())
	}
}
