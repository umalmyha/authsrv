package args

import (
	"os"
)

type ParsedArgs struct {
	args []string
}

func Parse() ParsedArgs {
	parsedArgs := ParsedArgs{args: make([]string, 0)}
	if len(os.Args) < 2 {
		return parsedArgs
	}
	parsedArgs.args = append(parsedArgs.args, os.Args[1:]...)
	return parsedArgs
}

func (a *ParsedArgs) At(index int) string {
	if index < 0 || index >= len(a.args) {
		return ""
	}
	return a.args[index]
}

func (a *ParsedArgs) KeyValueAt(index int) (string, string) {
	arg := a.At(index)
	if arg == "" {
		return "", ""
	}
	return KeyValue(arg)
}

func (a *ParsedArgs) Len() int {
	return len(a.args)
}

func (a *ParsedArgs) Iterator() *iterator {
	iter := iterator{
		args: make([]string, len(a.args)),
	}
	copy(iter.args, a.args)
	return &iter
}
