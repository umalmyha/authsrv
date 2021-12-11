package main

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/command"
)

func main() {
	if err := loadEnv(); err != nil {
		fmt.Printf("error occured on loading environment variables: %s", err.Error())
	}

	if err := run(); err != nil {
		fmt.Printf("error occurred during command execution: %s", err.Error())
	}
}

func run() error {
	args := args.Parse()

	var cmd command.Executor
	switch args.At(0) {
	case "createuser":
		cmd = command.NewCreateUserCommand(args)
	default:
		cmd = command.NewHelpCommand()
	}

	return cmd.Run()
}

func loadEnv() error {
	if os.Getenv("APP_ENV") != "production" { // TODO: add normal handling later
		return godotenv.Load()
	}
	return nil
}
