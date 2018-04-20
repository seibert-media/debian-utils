package command

type Command interface {
	Do() error
	Undo() error
}
