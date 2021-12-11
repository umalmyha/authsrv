package input

import (
	"fmt"
	"strings"
)

type OptionInputConfig struct {
	Prompt      string
	Default     string
	IsMandatory bool
	Options     []string
}

type optionInput struct {
	cfg OptionInputConfig
}

func NewOptionInput(cfg OptionInputConfig) *optionInput {
	if cfg.Options == nil {
		cfg.Options = make([]string, 0)
	}
	return &optionInput{
		cfg: cfg,
	}
}

func (i *optionInput) Read() (string, error) {
	input := NewSimpleInput(Config{
		Prompt:  i.cfg.Prompt,
		Default: string(i.cfg.Default),
	})

	var option string
	for {
		value, err := input.Read()
		if err != nil {
			return "", nil
		}

		if i.isInOptionsList(value) {
			option = value
			break
		}

		if value != "" {
			i.listOptions()
		}
	}

	return option, nil
}

func (i *optionInput) isInOptionsList(value string) bool {
	for _, option := range i.cfg.Options {
		if value == option {
			return true
		}
	}
	return false
}

func (i *optionInput) listOptions() {
	var optionsStr strings.Builder
	optionsStr.WriteString("abailable options are: ")

	length := len(i.cfg.Options)
	for i, option := range i.cfg.Options {
		optionsStr.WriteString(option)
		if !(i == length-1) {
			optionsStr.WriteString(", ")
		}
	}

	fmt.Println("incorrect input")
	fmt.Println(optionsStr.String())
}
