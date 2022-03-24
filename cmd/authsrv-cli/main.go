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
	case "createscope":
		cmd = command.NewCreateScopeCommand(args)
	case "createrole":
		cmd = command.NewCreateRoleCommand(args)
	case "assignscope":
		cmd = command.NewAssignScopeCommand(args)
	case "unassignscope":
		cmd = command.NewUnassignScopeCommand(args)
	case "assignrole":
		cmd = command.NewAssignRoleCommand(args)
	case "unassignrole":
		cmd = command.NewUnassignRoleCommand(args)
	case "genkeys":
		cmd = command.NewGenKeysCommand(args)
	default:
		cmd = command.NewHelpCommand()
	}

	return cmd.Run()
}
