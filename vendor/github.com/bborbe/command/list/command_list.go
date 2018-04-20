package list

import (
	"fmt"

	"log"

	"github.com/bborbe/command"
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
				if err := l.commands[i].Undo(); err != nil {
					log.Printf("undo command failed: %v\n", err)
				}
			}
			return fmt.Errorf("execute commands failed: %v", err)
		}
	}
	return nil
}
