package command

import (
	"context"
	"log"
	"time"

	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infra"
	"github.com/umalmyha/authsrv/internal/infra/service"
)

type assignScopeCommand struct {
	*LoggingCommand
	args args.ParsedArgs
}

type assignScopeCommandOptions struct {
	scope string
	role  string
	help  bool
}

func NewAssignScopeCommand(args args.ParsedArgs, logger *log.Logger) Executor {
	return &assignScopeCommand{
		LoggingCommand: &LoggingCommand{logger: logger},
		args:           args,
	}
}

func (c *assignScopeCommand) Run() error {
	options := c.extractOptions()
	if options.help {
		c.Help()
		return nil
	}

	var err error
	scopeName := options.scope
	if scopeName == "" {
		scopeName, err = input.NewSimpleInput(input.Config{Prompt: "scope", IsMandatory: true}).Read()
		if err != nil {
			return err
		}
	}

	roleName := options.role
	if roleName == "" {
		roleName, err = input.NewSimpleInput(input.Config{Prompt: "to role", IsMandatory: true}).Read()
		if err != nil {
			return err
		}
	}

	db, err := infra.ConnectToDb()
	if err != nil {
		return err
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.NewRoleService(db).AssignScope(ctx, roleName, scopeName); err != nil {
		return err
	}

	logger := c.Logger()
	logger.Printf("scope '%s' is assigned to role %s successfully", scopeName, roleName)
	logger.Println()

	return nil
}

func (c *assignScopeCommand) Help() {
	logger := c.Logger()
	logger.Println("assignscope - command assigns scope to role")
	logger.Println("options:")
	logger.Println("  --help - show help")
	logger.Println("  --scope - specify scope name")
	logger.Println("  --to - specify role name")
	logger.Println("example:")
	logger.Println("  assignscope --scope=scope1 --to=role1")
}

func (c *assignScopeCommand) extractOptions() assignScopeCommandOptions {
	options := assignScopeCommandOptions{}

	iter := c.args.Iterator()
	for iter.HasNext() {
		nextOpt := iter.Next()
		option, value := args.KeyValue(nextOpt)
		switch option {
		case "--help":
			options.help = true
		case "--scope":
			options.scope = value
		case "--to":
			options.role = value
		}
	}

	return options
}
