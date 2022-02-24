package command

import "fmt"

type helpCommand struct {
	execs []Executor
}

func NewHelpCommand() Executor {
	return &helpCommand{
		execs: []Executor{
			&createUserCommand{},
			&genKeysCommand{},
		},
	}
}

func (c *helpCommand) Run() error {
	c.Help()
	return nil
}

func (c *helpCommand) Help() {
	fmt.Println("--- Auth server CLI tool ---")
	fmt.Println("Available commands:")
	fmt.Println()
	for _, exec := range c.execs {
		exec.Help()
		fmt.Println()
	}
}
