package command_list

import (
	"fmt"

	"github.com/bborbe/debian/command"
)

type CommandList interface {
	Add(command command.Command)
	Run() error
}

type commandList struct {
	commands []command.Command
}

func New() *commandList {
	c := new(commandList)
	return c
}

func (l *commandList) Add(command command.Command) {
	l.commands = append(l.commands, command)
}

func (l *commandList) Run() error {
	for pos, c := range l.commands {
		err := c.Do()
		if err != nil {
			for i := pos; i >= 0; i-- {
				l.commands[i].Undo()
			}
			return fmt.Errorf("execute commands failed: %v", err)
		}
	}
	return nil
}
