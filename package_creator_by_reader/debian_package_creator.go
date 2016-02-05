package package_creator_by_reader

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/bborbe/command"
	command_adapter "github.com/bborbe/command/adapter"
	command_list "github.com/bborbe/command/list"
	debian_config "github.com/bborbe/debian_utils/config"
	debian_config_builder "github.com/bborbe/debian_utils/config_builder"
	debian_package_creator "github.com/bborbe/debian_utils/package_creator"
	"github.com/bborbe/log"
)

type Creator interface {
	CreatePackage(fileReader io.Reader, config *debian_config.Config, sourceDir string, targetDir string) error
}

var logger = log.DefaultLogger

type CommandListProvider func() command_list.CommandList
type ExtractFile func(fileReader io.Reader, targetDir string) error

type creator struct {
	commandListProvider CommandListProvider
	packageCreator      debian_package_creator.PackageCreator
	extractFile         ExtractFile
}

func New(commandListProvider CommandListProvider, debianPackageCreator debian_package_creator.PackageCreator, extractFile ExtractFile) *creator {
	d := new(creator)
	d.commandListProvider = commandListProvider
	d.packageCreator = debianPackageCreator
	d.extractFile = extractFile
	return d
}

func (d *creator) CreatePackage(fileReader io.Reader, config *debian_config.Config, sourceDir string, targetDir string) error {
	b := new(builder)
	b.packageCreator = d.packageCreator
	b.extractFile = d.extractFile
	b.commandList = d.commandListProvider()
	b.config = config
	b.fileReader = fileReader
	b.commandList.Add(b.createWorkingDirectoryCommand())
	b.commandList.Add(b.extractCommand())
	b.commandList.Add(b.createDebianPackageCommand())
	b.commandList.Add(b.removeWorkingDirectoryCommand())
	b.sourceDir = sourceDir
	b.targetDir = targetDir
	return b.commandList.Run()
}

type builder struct {
	extractFile      ExtractFile
	commandList      command_list.CommandList
	packageCreator   debian_package_creator.PackageCreator
	workingdirectory string
	fileReader       io.Reader
	config           *debian_config.Config
	path             string
	targetDir        string
	sourceDir        string
}

func (b *builder) createWorkingDirectoryCommand() command.Command {
	return command_adapter.New(func() error {
		logger.Debugf("create working directory")
		var err error
		if b.workingdirectory, err = ioutil.TempDir("", "create-debian-package"); err != nil {
			return err
		}
		logger.Debugf("working directory %s create", b.workingdirectory)
		return nil
	}, func() error {
		return os.RemoveAll(b.workingdirectory)
	})
}

func (b *builder) extractCommand() command.Command {
	return command_adapter.New(func() error {
		return b.extractFile(b.fileReader, b.workingdirectory)
	}, func() error { return nil })
}

func (b *builder) createDebianPackageCommand() command.Command {
	return command_adapter.New(func() error {
		configBuilder := debian_config_builder.NewWithConfig(b.config)
		configBuilder.AddFile(joinDirs(b.workingdirectory, b.sourceDir), b.targetDir)
		return b.packageCreator.CreatePackage(configBuilder.Build())
	}, func() error { return nil })
}

func (b *builder) removeWorkingDirectoryCommand() command.Command {
	return command_adapter.New(func() error {
		logger.Debugf("clean working directory %s", b.workingdirectory)
		if err := os.RemoveAll(b.workingdirectory); err != nil {
			return err
		}
		logger.Debugf("working directory %s cleaned", b.workingdirectory)
		return nil
	}, func() error {
		return nil
	})
}

func joinDirs(dirs ...string) string {
	result := ""
	first := true
	for _, dir := range dirs {
		if len(dir) > 0 {
			if first {
				first = false
				result = dir
			} else {
				result = fmt.Sprintf("%s%s%s", result, string(os.PathSeparator), dir)
			}
		}
	}
	return result
}
