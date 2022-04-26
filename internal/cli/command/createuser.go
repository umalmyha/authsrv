package command

import (
	"context"
	"log"
	"time"

	"github.com/umalmyha/authsrv/internal/business/user"
	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infra"
	"github.com/umalmyha/authsrv/internal/infra/service"
	dbredis "github.com/umalmyha/authsrv/pkg/database/redis"
)

type createUserCommand struct {
	*LoggingCommand
	args args.ParsedArgs
}

type createUserCommandOptions struct {
	help     bool
	username string
	password string
	isSuper  bool
}

func NewCreateUserCommand(args args.ParsedArgs, logger *log.Logger) Executor {
	return &createUserCommand{
		LoggingCommand: &LoggingCommand{logger: logger},
		args:           args,
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

	db, err := infra.ConnectToDb()
	if err != nil {
		return err
	}
	defer db.Close()

	redisOpts, err := infra.RedisOptions()
	if err != nil {
		return err
	}

	rdb, err := dbredis.Connect(redisOpts)
	if err != nil {
		return err
	}

	jwtCfg, err := infra.JwtConfig()
	if err != nil {
		return err
	}

	rfrCfg, err := infra.RefreshTokenConfig()
	if err != nil {
		return err
	}

	passCfg, err := infra.PasswordConfig()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	srv := service.NewAuthService(db, rdb, jwtCfg, rfrCfg, passCfg)
	nu := user.NewUserDto{
		Username:        username,
		Password:        password,
		ConfirmPassword: password,
		IsSuperuser:     false,
	}
	if err := srv.Signup(ctx, nu); err != nil {
		return err
	}

	logger := c.Logger()
	logger.Printf("user '%s' is created successfully", username)
	logger.Println()

	return nil
}

func (c *createUserCommand) Help() {
	logger := c.Logger()
	logger.Println("createuser - command creates new user")
	logger.Println("options:")
	logger.Println("  --help - show help")
	logger.Println("  --username - specify username")
	logger.Println("  --password - specify password")
	logger.Println("  --issuper - create superuser")
	logger.Println("example:")
	logger.Println("  createuser --usename=user1 --password=initial1 --issuper")
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
