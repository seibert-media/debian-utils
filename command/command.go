package command

type Command interface {
	Do() error
	Undo() error
}

type command struct {
	do   func() error
	undo func() error
}

func New(do func() error, undo func() error) *command {
	c := new(command)
	c.do = do
	c.undo = undo
	return c
}

func (c *command) Do() error {
	return c.do()
}

func (c *command) Undo() error {
	return c.undo()
}
