package command

import (
	"context"
	"fmt"
	"time"

	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infrastruct"
	"github.com/umalmyha/authsrv/internal/service"
)

type unassignScopeCommand struct {
	args args.ParsedArgs
}

type unassignScopeCommandOptions struct {
	scope string
	role  string
	help  bool
}

func NewUnassignScopeCommand(args args.ParsedArgs) Executor {
	return &unassignScopeCommand{
		args: args,
	}
}

func (c *unassignScopeCommand) Run() error {
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
		roleName, err = input.NewSimpleInput(input.Config{Prompt: "from role", IsMandatory: true}).Read()
		if err != nil {
			return err
		}
	}

	db, err := infrastruct.ConnectToDb()
	if err != nil {
		return err
	}
	defer db.Close()

	logger, err := infrastruct.NewCliZapLogger()
	if err != nil {
		return err
	}
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.NewRoleService(db).AssignScope(ctx, roleName, scopeName); err != nil {
		return err
	}

	fmt.Printf("scope '%s' is unassigned from role %s successfully", scopeName, roleName)
	fmt.Println()

	return nil
}

func (c *unassignScopeCommand) Help() {
	fmt.Println("unassignscope - command unassigns scope from role")
	fmt.Println("options:")
	fmt.Println("  --help - show help")
	fmt.Println("  --scope - specify scope name")
	fmt.Println("  --from - specify role name")
	fmt.Println("example:")
	fmt.Println("  unassignscope --scope=scope1 --from=role1")
}

func (c *unassignScopeCommand) extractOptions() unassignScopeCommandOptions {
	options := unassignScopeCommandOptions{}

	iter := c.args.Iterator()
	for iter.HasNext() {
		nextOpt := iter.Next()
		option, value := args.KeyValue(nextOpt)
		switch option {
		case "--help":
			options.help = true
		case "--scope":
			options.scope = value
		case "--from":
			options.role = value
		}
	}

	return options
}
