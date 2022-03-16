package main

import (
	"fmt"

	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/command"
	"github.com/umalmyha/authsrv/internal/infrastruct"
)

func main() {
	if err := infrastruct.LoadEnv(); err != nil {
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
	case "genkeys":
		cmd = command.NewGenKeysCommand(args)
	default:
		cmd = command.NewHelpCommand()
	}

	return cmd.Run()
}
