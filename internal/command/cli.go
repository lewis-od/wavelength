package command

type Command interface {
	Run(arguments []string)
}
