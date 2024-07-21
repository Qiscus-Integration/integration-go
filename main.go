package main

import (
	"integration-go/cmd"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	cmd.Execute()
}
