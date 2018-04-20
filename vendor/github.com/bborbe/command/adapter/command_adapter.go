package adapter

type commandAdapter struct {
	do   func() error
	undo func() error
}

func New(do func() error, undo func() error) *commandAdapter {
	c := new(commandAdapter)
	c.do = do
	c.undo = undo
	return c
}

func (c *commandAdapter) Do() error {
	return c.do()
}

func (c *commandAdapter) Undo() error {
	return c.undo()
}
