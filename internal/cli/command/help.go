package command

import (
	"log"
)

type helpCommand struct {
	*LoggingCommand
	execs []Executor
}

func NewHelpCommand(logger *log.Logger) Executor {
	return &helpCommand{
		LoggingCommand: &LoggingCommand{logger: logger},
		execs: []Executor{
			&createUserCommand{},
			&createScopeCommand{},
			&createRoleCommand{},
			&assignScopeCommand{},
			&unassignScopeCommand{},
			&assignRoleCommand{},
			&unassignRoleCommand{},
			&genKeysCommand{},
		},
	}
}

func (c *helpCommand) Run() error {
	c.Help()
	return nil
}

func (c *helpCommand) Help() {
	logger := c.Logger()
	logger.Println("--- Auth server CLI tool ---")
	logger.Println("Available commands:")
	logger.Println()
	for _, exec := range c.execs {
		exec.Help()
		logger.Println()
	}
}
