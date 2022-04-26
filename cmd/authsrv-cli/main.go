package main

import (
	"fmt"

	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/command"
	"github.com/umalmyha/authsrv/internal/infra"
)

func main() {
	if err := infra.LoadEnv(); err != nil {
		fmt.Printf("error occured on loading environment variables: %s", err.Error())
	}

	if err := run(); err != nil {
		fmt.Printf("error occurred during command execution: %s", err.Error())
	}
}

func run() error {
	args := args.Parse()
	logger, err := infra.NewCliZapLogger()
	if err != nil {
		return err
	}

	var cmd command.Executor
	switch args.At(0) {
	case "createuser":
		cmd = command.NewCreateUserCommand(args, logger)
	case "createscope":
		cmd = command.NewCreateScopeCommand(args, logger)
	case "createrole":
		cmd = command.NewCreateRoleCommand(args, logger)
	case "assignscope":
		cmd = command.NewAssignScopeCommand(args, logger)
	case "unassignscope":
		cmd = command.NewUnassignScopeCommand(args, logger)
	case "assignrole":
		cmd = command.NewAssignRoleCommand(args, logger)
	case "unassignrole":
		cmd = command.NewUnassignRoleCommand(args, logger)
	case "genkeys":
		cmd = command.NewGenKeysCommand(args, logger)
	default:
		cmd = command.NewHelpCommand(logger)
	}

	return cmd.Run()
}
