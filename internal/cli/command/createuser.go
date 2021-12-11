package command

import (
	"context"
	"fmt"

	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/user"
	"github.com/umalmyha/authsrv/internal/user/dto"
	"github.com/umalmyha/authsrv/internal/user/store"
)

type createUserCommandOptions struct {
	help     bool
	username string
	password string
	isSuper  bool
}

type createUserCommand struct {
	args args.ParsedArgs
}

func NewCreateUserCommand(args args.ParsedArgs) Executor {
	return &createUserCommand{
		args: args,
	}
}

func (c *createUserCommand) Run() error {
	options := c.extractOptions()

	if options.help {
		c.Help()
		return nil
	}

	var err error
	username := options.username
	if username == "" {
		username, err = input.NewSimpleInput(input.Config{Prompt: "username", IsMandatory: true}).Read()
		if err != nil {
			return err
		}
	}

	password := options.password
	if password == "" {
		password, err = input.NewPasswordInput().Read()
		if err != nil {
			return err
		}
	}

	// TODO: Improve creation process later
	db, err := connectToDb()
	if err != nil {
		return err
	}

	logger, err := newZapLogger()
	if err != nil {
		return err
	}

	srv := user.Service(logger, store.NewStore(db))
	nu := dto.NewUser{
		Username:        username,
		Password:        password,
		ConfirmPassword: password,
		IsSuperuser:     true,
	}
	if _, err := srv.CreateUser(context.Background(), nu); err != nil {
		return err
	}

	fmt.Printf("user %s is created successfully", username)
	fmt.Println()

	return nil
}

func (c *createUserCommand) Help() {
	fmt.Println("createuser - command creates new user")
	fmt.Println("options:")
	fmt.Println("  --help - show help")
	fmt.Println("  --username - specify username")
	fmt.Println("  --password - specify password")
	fmt.Println("  --issuper - create superuser")
}

func (c *createUserCommand) extractOptions() createUserCommandOptions {
	options := createUserCommandOptions{}

	iter := c.args.Iterator()
	for iter.HasNext() {
		nextOpt := iter.Next()
		option, value := args.KeyValue(nextOpt)
		switch option {
		case "--help":
			options.help = true
		case "--username":
			options.username = value
		case "--password":
			options.password = value
		case "--issuper":
			options.isSuper = true
		}
	}

	return options
}
