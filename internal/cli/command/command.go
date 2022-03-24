package command

type Executor interface {
	Run() error
	Help()
}
