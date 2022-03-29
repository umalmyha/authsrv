package command

import (
	"context"
	"fmt"
	"time"

	"github.com/umalmyha/authsrv/internal/business/scope"
	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infra"
	"github.com/umalmyha/authsrv/internal/service"
)

type createScopeCommand struct {
	args args.ParsedArgs
}

type createScopeCommandOptions struct {
	name string
	help bool
}

func NewCreateScopeCommand(args args.ParsedArgs) Executor {
	return &createScopeCommand{
		args: args,
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

	logger, err := infra.NewCliZapLogger()
	if err != nil {
		return err
	}
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv := service.NewScopeService(db)
	if err := srv.CreateScope(ctx, scope.NewScopeDto{Name: name}); err != nil {
		return err
	}

	fmt.Printf("scope '%s' is created successfully", name)
	fmt.Println()

	return nil
}

func (c *createScopeCommand) Help() {
	fmt.Println("createscope - command creates new scope")
	fmt.Println("options:")
	fmt.Println("  --help - show help")
	fmt.Println("  --name - specify scope name")
	fmt.Println("example:")
	fmt.Println("  createscope --name=scope1")
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
