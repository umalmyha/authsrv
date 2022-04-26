package command

import "log"

type Executor interface {
	Logger() *log.Logger
	Run() error
	Help()
}

type LoggingCommand struct {
	logger *log.Logger
}

func (cmd *LoggingCommand) Logger() *log.Logger {
	return cmd.logger
}
