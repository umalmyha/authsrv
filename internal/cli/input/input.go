package input

import (
	"fmt"
	"strings"
)

type Config struct {
	Prompt      string
	Default     string
	IsMandatory bool
}

func okValue(input string, def string, mand bool) (bool, string) {
	if input != "" {
		return true, input
	}

	input = def
	if mand && input == "" {
		return false, input
	}

	return true, input
}

func promptStr(label string, def string, mand bool) string {
	var prompt strings.Builder

	prompt.WriteString(label)
	if def != "" {
		prompt.WriteString(fmt.Sprintf("(default %s)", def))
	} else if !mand {
		prompt.WriteString("(optional)")
	}

	prompt.WriteString(": ")
	return prompt.String()
}
