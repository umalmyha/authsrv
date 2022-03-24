package command

import (
	"context"
	"fmt"
	"time"

	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infrastruct"
	"github.com/umalmyha/authsrv/internal/service"
	dbredis "github.com/umalmyha/authsrv/pkg/database/redis"
)

type assignRoleCommand struct {
	args args.ParsedArgs
}

type assignRoleCommandOptions struct {
	role     string
	username string
	help     bool
}

func NewAssignRoleCommand(args args.ParsedArgs) Executor {
	return &assignRoleCommand{
		args: args,
	}
}

func (c *assignRoleCommand) Run() error {
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
		username, err = input.NewSimpleInput(input.Config{Prompt: "to user", IsMandatory: true}).Read()
		if err != nil {
			return err
		}
	}

	db, err := infrastruct.ConnectToDb()
	if err != nil {
		return err
	}
	defer db.Close()

	redisOpts, err := infrastruct.RedisOptions()
	if err != nil {
		return err
	}

	rdb, err := dbredis.Connect(redisOpts)
	if err != nil {
		return err
	}

	logger, err := infrastruct.NewCliZapLogger()
	if err != nil {
		return err
	}
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.NewUserService(db, rdb).AssignRole(ctx, username, roleName); err != nil {
		return err
	}

	fmt.Printf("role '%s' is assigned to user %s successfully", roleName, username)
	fmt.Println()

	return nil
}

func (c *assignRoleCommand) Help() {
	fmt.Println("assignrole - command assigns role to user")
	fmt.Println("options:")
	fmt.Println("  --help - show help")
	fmt.Println("  --role - specify role name")
	fmt.Println("  --to - specify username")
	fmt.Println("example:")
	fmt.Println("  assignrole --role=role1 --to=username1")
}

func (c *assignRoleCommand) extractOptions() assignRoleCommandOptions {
	options := assignRoleCommandOptions{}

	iter := c.args.Iterator()
	for iter.HasNext() {
		nextOpt := iter.Next()
		option, value := args.KeyValue(nextOpt)
		switch option {
		case "--help":
			options.help = true
		case "--role":
			options.role = value
		case "--to":
			options.username = value
		}
	}

	return options
}
