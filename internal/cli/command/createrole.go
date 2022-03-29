package command

import (
	"context"
	"fmt"
	"time"

	"github.com/umalmyha/authsrv/internal/business/role"
	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infra"
	"github.com/umalmyha/authsrv/internal/service"
)

type createRoleCommand struct {
	args args.ParsedArgs
}

type createRoleCommandOptions struct {
	name string
	help bool
}

func NewCreateRoleCommand(args args.ParsedArgs) Executor {
	return &createRoleCommand{
		args: args,
	}
}

func (c *createRoleCommand) Run() error {
	options := c.extractOptions()
	if options.help {
		c.Help()
		return nil
	}

	var err error
	name := options.name
	if name == "" {
		name, err = input.NewSimpleInput(input.Config{Prompt: "name", IsMandatory: true}).Read()
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

	srv := service.NewRoleService(db)
	if err := srv.CreateRole(ctx, role.NewRoleDto{Name: name}); err != nil {
		return err
	}

	fmt.Printf("role '%s' is created successfully", name)
	fmt.Println()

	return nil
}

func (c *createRoleCommand) Help() {
	fmt.Println("createrole - command creates new role")
	fmt.Println("options:")
	fmt.Println("  --help - show help")
	fmt.Println("  --name - specify role name")
	fmt.Println("example:")
	fmt.Println("  createrole --name=role1")
}

func (c *createRoleCommand) extractOptions() createRoleCommandOptions {
	options := createRoleCommandOptions{}

	iter := c.args.Iterator()
	for iter.HasNext() {
		nextOpt := iter.Next()
		option, value := args.KeyValue(nextOpt)
		switch option {
		case "--help":
			options.help = true
		case "--name":
			options.name = value
		}
	}

	return options
}
