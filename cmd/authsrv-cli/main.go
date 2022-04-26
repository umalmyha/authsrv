package main

import (
	"log"

	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/command"
	"github.com/umalmyha/authsrv/internal/infra"
	"go.uber.org/zap"
)

func main() {
	logger, err := infra.NewCliZapLogger()
	if err != nil {
		log.Fatal(err)
	}
	defer logger.Sync()

	stdlogger := zap.NewStdLog(logger.Desugar())
	if err := infra.LoadEnv(); err != nil {
		stdlogger.Printf("error occured on loading environment variables: %s", err.Error())
	}

	if err := run(stdlogger); err != nil {
		stdlogger.Printf("error occurred during command execution: %s", err.Error())
	}
}

func run(logger *log.Logger) error {
	args := args.Parse()

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
