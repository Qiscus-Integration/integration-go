package main

import (
	"integration-go/cmd"

	_ "github.com/joho/godotenv/autoload"
	"github.com/rs/zerolog/log"
)

func main() {
	if err := cmd.NewRootCommand().Execute(); err != nil {
		log.Fatal().Msgf("failed run app: %s", err.Error())
	}
}
