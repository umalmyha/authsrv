package command

import (
	"context"
	"log"
	"time"

	"github.com/umalmyha/authsrv/internal/business/scope"
	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infra"
	"github.com/umalmyha/authsrv/internal/infra/service"
)

type createScopeCommand struct {
	*LoggingCommand
	args args.ParsedArgs
}

type createScopeCommandOptions struct {
	name string
	help bool
}

func NewCreateScopeCommand(args args.ParsedArgs, logger *log.Logger) Executor {
	return &createScopeCommand{
		LoggingCommand: &LoggingCommand{logger: logger},
		args:           args,
	}
}

func (c *createScopeCommand) Run() error {
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

	srv := service.NewScopeService(db)
	if err := srv.CreateScope(ctx, scope.NewScopeDto{Name: name}); err != nil {
		return err
	}

	logger := c.Logger()
	logger.Printf("scope '%s' is created successfully", name)
	logger.Println()

	return nil
}

func (c *createScopeCommand) Help() {
	logger := c.Logger()
	logger.Println("createscope - command creates new scope")
	logger.Println("options:")
	logger.Println("  --help - show help")
	logger.Println("  --name - specify scope name")
	logger.Println("example:")
	logger.Println("  createscope --name=scope1")
}

func (c *createScopeCommand) extractOptions() createScopeCommandOptions {
	options := createScopeCommandOptions{}

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
