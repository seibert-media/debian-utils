package package_creator
import (
	"github.com/bborbe/debian/command"
	"github.com/bborbe/debian/command_list"
)

type Builder interface {
	Build() error
}

type builder struct {
	name        string
	version     string
	commandList command_list.CommandList
}

func New() *builder {
	return new (builder)
}

func (b *builder) Name(name string) *builder {
	b.name = name
	return b
}

func (b *builder) Version(version string) *builder {
	b.version = version
	return b
}

func (b *builder) Build() error {
	b.commandList.Add(b.createWorkingDirectoryCommand())
	b.commandList.Add(b.createDebianFolderCommand())
	b.commandList.Add(b.createDebianControlCommand())
	b.commandList.Add(b.createDebianPackageCommand())
	b.commandList.Add(b.cleanWorkingDirectoryCommand())
	return b.commandList.Run()
}

func (b *builder) createWorkingDirectoryCommand() command.Command {
	return command.New(func() error {
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) createDebianFolderCommand() command.Command {
	return command.New(func() error {
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) createDebianControlCommand() command.Command {
	return command.New(func() error {
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) createDebianPackageCommand() command.Command {
	return command.New(func() error {
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) cleanWorkingDirectoryCommand() command.Command {
	return command.New(func() error {
		return nil
	}, func() error {
		return nil
	})
}
