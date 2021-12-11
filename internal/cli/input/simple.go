package input

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type simpleInput struct {
	cfg Config
}

func NewSimpleInput(cfg Config) *simpleInput {
	return &simpleInput{
		cfg: cfg,
	}
}

func (i *simpleInput) Read() (string, error) {
	scanner := bufio.NewScanner(os.Stdin)

	var input string
	for {
		fmt.Print(promptStr(i.cfg.Prompt, i.cfg.Default, i.cfg.IsMandatory))

		scanner.Scan()
		if err := scanner.Err(); err != nil {
			return "", err
		}
		fmt.Println()

		if ok, value := okValue(scanner.Text(), i.cfg.Default, i.cfg.IsMandatory); ok {
			input = value
			break
		}
	}

	return strings.TrimSpace(input), nil
}
