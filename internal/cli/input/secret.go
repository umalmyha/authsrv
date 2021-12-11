package input

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type secretInput struct {
	cfg Config
}

func NewSecretInput(cfg Config) *secretInput {
	return &secretInput{
		cfg: cfg,
	}
}

func (i *secretInput) Read() (string, error) {
	fd := int(os.Stdin.Fd())

	initState, err := term.MakeRaw(fd)
	if err != nil {
		return "", err
	}
	defer term.Restore(fd, initState)

	var secret string
	for {
		fmt.Print(promptStr(i.cfg.Prompt, i.cfg.Default, i.cfg.IsMandatory))

		bytePass, err := term.ReadPassword(fd)
		if err != nil {
			return "", nil
		}
		fmt.Println()

		if ok, value := okValue(string(bytePass), i.cfg.Default, i.cfg.IsMandatory); ok {
			secret = value
			break
		}
	}

	return secret, nil
}
