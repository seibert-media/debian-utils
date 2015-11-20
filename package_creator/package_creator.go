package package_creator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/bborbe/log"
	debian_command "github.com/bborbe/debian/command"
	debian_command_list "github.com/bborbe/debian/command_list"
	debian_config "github.com/bborbe/debian/config"
	debian_copier "github.com/bborbe/debian/copier"
)

type PackageCreator interface {
	CreatePackage(config *debian_config.Config) error
}

type packageCreator struct {
	commandListProvider CommandListProvider
	copier              debian_copier.Copier
}

type builder struct {
	config           *debian_config.Config
	command_list     debian_command_list.CommandList
	copier           debian_copier.Copier
	workingdirectory string
}

var logger = log.DefaultLogger

type CommandListProvider func() debian_command_list.CommandList

func New(commandListProvider CommandListProvider, copier debian_copier.Copier) *packageCreator {
	p := new(packageCreator)
	p.commandListProvider = commandListProvider
	p.copier = copier
	return p
}

func (p *packageCreator) CreatePackage(config *debian_config.Config) error {
	b := new(builder)
	b.command_list = p.commandListProvider()
	b.copier = p.copier
	b.config = config
	logger.Debug("Build")
	b.command_list.Add(b.validateCommand())
	b.command_list.Add(b.createWorkingDirectoryCommand())
	b.command_list.Add(b.createDebianFolderCommand())
	b.command_list.Add(b.createDebianControlCommand())
	b.command_list.Add(b.copyFilesToWorkingDirectoryCommand())
	b.command_list.Add(b.createDebianPackageCommand())
	b.command_list.Add(b.copyDebianPackageCommand())
	b.command_list.Add(b.cleanWorkingDirectoryCommand())
	return b.command_list.Run()
}

func (b *builder) validateCommand() debian_command.Command {
	return debian_command.New(func() error {
		logger.Debug("validate")
		if len(b.config.Files) == 0 {
			return fmt.Errorf("add at least one file")
		}
		if len(b.config.Name) == 0 {
			return fmt.Errorf("name missing")
		}
		if len(b.config.Version) == 0 {
			return fmt.Errorf("version missing")
		}
		logger.Debug("validate success")
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) createWorkingDirectoryCommand() debian_command.Command {
	return debian_command.New(func() error {
		logger.Debugf("create working directory")
		var err error
		if b.workingdirectory, err = ioutil.TempDir("", fmt.Sprintf("%s_%s", b.config.Name, b.config.Version)); err != nil {
			return err
		}
		logger.Debugf("working directory %s create", b.workingdirectory)
		return nil
	}, func() error {
		return os.RemoveAll(b.workingdirectory)
	})
}

func (b *builder) createDebianFolderCommand() debian_command.Command {
	return debian_command.New(func() error {
		path := fmt.Sprintf("%s/%s_%s/DEBIAN", b.workingdirectory, b.config.Name, b.config.Version)
		logger.Debugf("create debian folder %s", path)
		if err := createDirectory(path); err != nil {
			return err
		}
		logger.Debugf("debian folder created")
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) createDebianControlCommand() debian_command.Command {
	return debian_command.New(func() error {
		logger.Debugf("write debian control")
		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s_%s/DEBIAN/control", b.workingdirectory, b.config.Name, b.config.Version), b.controlContent(), 0644); err != nil {
			return err
		}
		logger.Debugf("debian control written")
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) controlContent() []byte {
	buffer := bytes.NewBufferString("")
	fmt.Fprintf(buffer, "Package: %s\n", b.config.Name)
	fmt.Fprintf(buffer, "Version: %s\n", b.config.Version)
	fmt.Fprintf(buffer, "Section: %s\n", b.config.Section)
	fmt.Fprintf(buffer, "Priority: %s\n", b.config.Priority)
	fmt.Fprintf(buffer, "Architecture: %s\n", b.config.Architecture)
	fmt.Fprintf(buffer, "Maintainer: %s\n", b.config.Maintainer)
	fmt.Fprintf(buffer, "Description: %s\n", b.config.Description)
	return buffer.Bytes()
}

func (b *builder) createDebianPackageCommand() debian_command.Command {
	return debian_command.New(func() error {
		logger.Debugf("create debian package")
		cmd := exec.Command("dpkg-deb", "--build", fmt.Sprintf("%s_%s", b.config.Name, b.config.Version))
		cmd.Dir = b.workingdirectory
		cmd.Stderr = os.Stderr
		//cmd.Stdout = os.Stdout
		if err := cmd.Run(); err != nil {
			return err
		}
		logger.Debugf("debian package created")
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) copyDebianPackageCommand() debian_command.Command {
	return debian_command.New(func() error {
		logger.Debugf("copy debian package")
		var dir string
		var err error
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		source := fmt.Sprintf("%s/%s_%s.deb", b.workingdirectory, b.config.Name, b.config.Version)
		target := fmt.Sprintf("%s/%s_%s.deb", dir, b.config.Name, b.config.Version)
		if err = b.copier.Copy(source, target); err != nil {
			return err
		}
		logger.Debugf("debian package copied")
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) copyFilesToWorkingDirectoryCommand() debian_command.Command {
	return debian_command.New(func() error {
		logger.Debugf("copy files")
		for _, file := range b.config.Files {
			var err error
			var directory string
			var filename string
			filename = fmt.Sprintf("%s/%s_%s%s", b.workingdirectory, b.config.Name, b.config.Version, file.Target)
			if directory, err = dirOf(filename); err != nil {
				return err
			}
			if err = createDirectory(directory); err != nil {
				return err
			}
			if err = b.copier.Copy(file.Source, filename); err != nil {
				return err
			}
		}
		logger.Debugf("all files copied")
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) cleanWorkingDirectoryCommand() debian_command.Command {
	return debian_command.New(func() error {
		logger.Debugf("clean working directory")
		if err := os.RemoveAll(b.workingdirectory); err != nil {
			return err
		}
		logger.Debugf("working directory cleaned")
		return nil
	}, func() error {
		return nil
	})
}

func createDirectory(directory string) error {
	logger.Debugf("create directory %s", directory)
	return os.MkdirAll(directory, 0755)
}

func dirOf(filename string) (string, error) {
	pos := strings.LastIndex(filename, "/")
	if pos != -1 {
		return filename[:pos + 1], nil
	}
	return "", fmt.Errorf("can't determine directory of file %s", filename)
}

