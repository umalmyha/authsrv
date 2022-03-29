package command

import (
	"context"
	"fmt"
	"time"

	"github.com/umalmyha/authsrv/internal/cli/args"
	"github.com/umalmyha/authsrv/internal/cli/input"
	"github.com/umalmyha/authsrv/internal/infra"
	"github.com/umalmyha/authsrv/internal/service"
	dbredis "github.com/umalmyha/authsrv/pkg/database/redis"
)

type unassignRoleCommand struct {
	args args.ParsedArgs
}

type unassignRoleCommandOptions struct {
	role     string
	username string
	help     bool
}

func NewUnassignRoleCommand(args args.ParsedArgs) Executor {
	return &unassignRoleCommand{
		args: args,
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

	logger, err := infra.NewCliZapLogger()
	if err != nil {
		return err
	}
	defer logger.Sync()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := service.NewUserService(db, rdb).UnassignRole(ctx, username, roleName); err != nil {
		return err
	}

	fmt.Printf("role '%s' is unassigned from user %s successfully", roleName, username)
	fmt.Println()

	return nil
}

func (c *unassignRoleCommand) Help() {
	fmt.Println("unassignrole - command unassigns role from user")
	fmt.Println("options:")
	fmt.Println("  --help - show help")
	fmt.Println("  --role - specify role name")
	fmt.Println("  --from - specify username")
	fmt.Println("example:")
	fmt.Println("  assignrole --role=role1 --from=username1")
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
