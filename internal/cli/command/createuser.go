package command

import (
	"context"
	"fmt"

	"github.com/umalmyha/authsrv/internal/business/user"
	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infrastruct"
	"github.com/umalmyha/authsrv/internal/service"
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

	jwtCfg, err := infrastruct.JwtConfig()
	if err != nil {
		return err
	}

	rfrCfg, err := infrastruct.RefreshTokenConfig()
	if err != nil {
		return err
	}

	passCfg, err := infrastruct.PasswordConfig()
	if err != nil {
		return err
	}

	srv := service.NewAuthService(db, jwtCfg, rfrCfg, passCfg)
	nu := user.NewUserDto{
		Username:        username,
		Password:        password,
		ConfirmPassword: password,
		IsSuperuser:     false,
	}
	if err := srv.Signup(context.Background(), nu); err != nil {
		return err
	}

	fmt.Printf("user '%s' is created successfully", username)
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
	fmt.Println("example:")
	fmt.Println("    createuser --usename=user1 --password=initial1 --issuper")
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
