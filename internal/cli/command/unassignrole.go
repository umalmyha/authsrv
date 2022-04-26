package command

import (
	"context"
	"log"
	"time"

	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infra"
	"github.com/umalmyha/authsrv/internal/infra/service"
	dbredis "github.com/umalmyha/authsrv/pkg/database/redis"
)

type unassignRoleCommand struct {
	*LoggingCommand
	args args.ParsedArgs
}

type unassignRoleCommandOptions struct {
	role     string
	username string
	help     bool
}

func NewUnassignRoleCommand(args args.ParsedArgs, logger *log.Logger) Executor {
	return &unassignRoleCommand{
		LoggingCommand: &LoggingCommand{logger: logger},
		args:           args,
	}
}

func (c *unassignRoleCommand) Run() error {
	options := c.extractOptions()
	if options.help {
		c.Help()
		return nil
	}

	var err error
	roleName := options.role
	if roleName == "" {
		roleName, err = input.NewSimpleInput(input.Config{Prompt: "role", IsMandatory: true}).Read()
		if err != nil {
			return err
		}
	}

	username := options.username
	if username == "" {
		username, err = input.NewSimpleInput(input.Config{Prompt: "from user", IsMandatory: true}).Read()
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

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.NewUserService(db, rdb).UnassignRole(ctx, username, roleName); err != nil {
		return err
	}

	logger := c.Logger()
	logger.Printf("role '%s' is unassigned from user %s successfully", roleName, username)
	logger.Println()

	return nil
}

func (c *unassignRoleCommand) Help() {
	logger := c.Logger()
	logger.Println("unassignrole - command unassigns role from user")
	logger.Println("options:")
	logger.Println("  --help - show help")
	logger.Println("  --role - specify role name")
	logger.Println("  --from - specify username")
	logger.Println("example:")
	logger.Println("  assignrole --role=role1 --from=username1")
}

func (c *unassignRoleCommand) extractOptions() unassignRoleCommandOptions {
	options := unassignRoleCommandOptions{}

	iter := c.args.Iterator()
	for iter.HasNext() {
		nextOpt := iter.Next()
		option, value := args.KeyValue(nextOpt)
		switch option {
		case "--help":
			options.help = true
		case "--role":
			options.role = value
		case "--from":
			options.username = value
		}
	}

	return options
}
