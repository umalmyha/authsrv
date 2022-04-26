package command

import (
	"context"
	"log"
	"time"

	"github.com/umalmyha/authsrv/internal/business/role"
	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infra"
	"github.com/umalmyha/authsrv/internal/infra/service"
)

type createRoleCommand struct {
	*LoggingCommand
	args args.ParsedArgs
}

type createRoleCommandOptions struct {
	name string
	help bool
}

func NewCreateRoleCommand(args args.ParsedArgs, logger *log.Logger) Executor {
	return &createRoleCommand{
		LoggingCommand: &LoggingCommand{logger: logger},
		args:           args,
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv := service.NewRoleService(db)
	if err := srv.CreateRole(ctx, role.NewRoleDto{Name: name}); err != nil {
		return err
	}

	logger := c.Logger()
	logger.Printf("role '%s' is created successfully", name)
	logger.Println()

	return nil
}

func (c *createRoleCommand) Help() {
	logger := c.Logger()
	logger.Println("createrole - command creates new role")
	logger.Println("options:")
	logger.Println("  --help - show help")
	logger.Println("  --name - specify role name")
	logger.Println("example:")
	logger.Println("  createrole --name=role1")
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
