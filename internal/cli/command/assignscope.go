package command

import (
	"context"
	"fmt"
	"time"

	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infra"
	"github.com/umalmyha/authsrv/internal/service"
)

type assignScopeCommand struct {
	args args.ParsedArgs
}

type assignScopeCommandOptions struct {
	scope string
	role  string
	help  bool
}

func NewAssignScopeCommand(args args.ParsedArgs) Executor {
	return &assignScopeCommand{
		args: args,
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

	logger, err := infra.NewCliZapLogger()
	if err != nil {
		return err
	}
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.NewRoleService(db).AssignScope(ctx, roleName, scopeName); err != nil {
		return err
	}

	fmt.Printf("scope '%s' is assigned to role %s successfully", scopeName, roleName)
	fmt.Println()

	return nil
}

func (c *assignScopeCommand) Help() {
	fmt.Println("assignscope - command assigns scope to role")
	fmt.Println("options:")
	fmt.Println("  --help - show help")
	fmt.Println("  --scope - specify scope name")
	fmt.Println("  --to - specify role name")
	fmt.Println("example:")
	fmt.Println("  assignscope --scope=scope1 --to=role1")
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
