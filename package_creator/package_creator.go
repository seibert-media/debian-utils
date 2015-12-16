package package_creator

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	command "github.com/bborbe/command"
	command_adapter "github.com/bborbe/command/adapter"
	command_list "github.com/bborbe/command/list"
	debian_config "github.com/bborbe/debian_utils/config"
	debian_copier "github.com/bborbe/debian_utils/copier"
	"github.com/bborbe/log"
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
	command_list     command_list.CommandList
	copier           debian_copier.Copier
	workingdirectory string
}

var logger = log.DefaultLogger

type CommandListProvider func() command_list.CommandList

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
	logger.Debug("CreatePackage")
	b.command_list.Add(b.validateCommand())
	b.command_list.Add(b.createWorkingDirectoryCommand())
	b.command_list.Add(b.createDebianFolderCommand())
	b.command_list.Add(b.createDebianControlCommand())
	b.command_list.Add(b.createDebianConffilesCommand())
	b.command_list.Add(b.createDebianCopyPrePostFilesCommand())
	b.command_list.Add(b.copyFilesToWorkingDirectoryCommand())
	b.command_list.Add(b.createDebianPackageCommand())
	b.command_list.Add(b.copyDebianPackageCommand())
	b.command_list.Add(b.cleanWorkingDirectoryCommand())
	return b.command_list.Run()
}

func (b *builder) validateCommand() command.Command {
	return command_adapter.New(func() error {
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

func (b *builder) createWorkingDirectoryCommand() command.Command {
	return command_adapter.New(func() error {
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

func (b *builder) createDebianCopyPrePostFilesCommand() command.Command {
	return command_adapter.New(func() error {
		logger.Debugf("copy pre post inst rm files")
		var err error
		path := fmt.Sprintf("%s/%s_%s/DEBIAN", b.workingdirectory, b.config.Name, b.config.Version)
		if len(b.config.Postinst) > 0 {
			if err = b.copier.Copy(b.config.Postinst, fmt.Sprintf("%s/postinst", path)); err != nil {
				return err
			}
		}
		if len(b.config.Postrm) > 0 {
			if err = b.copier.Copy(b.config.Postrm, fmt.Sprintf("%s/postrm", path)); err != nil {
				return err
			}
		}
		if len(b.config.Preinst) > 0 {
			if err = b.copier.Copy(b.config.Preinst, fmt.Sprintf("%s/preinst", path)); err != nil {
				return err
			}
		}
		if len(b.config.Prerm) > 0 {
			if err = b.copier.Copy(b.config.Prerm, fmt.Sprintf("%s/prerm", path)); err != nil {
				return err
			}
		}
		logger.Debugf("pre post inst rm files copied")
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) createDebianFolderCommand() command.Command {
	return command_adapter.New(func() error {
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

func (b *builder) createDebianControlCommand() command.Command {
	return command_adapter.New(func() error {
		logger.Debugf("write debian control")
		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s_%s/DEBIAN/control", b.workingdirectory, b.config.Name, b.config.Version), controlContent(*b.config), 0644); err != nil {
			return err
		}
		logger.Debugf("debian control written")
		return nil
	}, func() error {
		return nil
	})
}

func controlContent(config debian_config.Config) []byte {
	buffer := bytes.NewBufferString("")
	fmt.Fprintf(buffer, "Package: %s\n", config.Name)
	fmt.Fprintf(buffer, "Version: %s\n", config.Version)
	fmt.Fprintf(buffer, "Section: %s\n", config.Section)
	fmt.Fprintf(buffer, "Priority: %s\n", config.Priority)
	fmt.Fprintf(buffer, "Architecture: %s\n", config.Architecture)
	fmt.Fprintf(buffer, "Maintainer: %s\n", config.Maintainer)
	fmt.Fprintf(buffer, "Description: %s\n", config.Description)
	if len(config.Depends) > 0 {
		fmt.Fprintf(buffer, "Depends: %s\n", strings.Join(config.Depends, ","))
	}
	return buffer.Bytes()
}

func (b *builder) createDebianConffilesCommand() command.Command {
	return command_adapter.New(func() error {
		logger.Debugf("write debian conffiles")
		content := conffilesContent(b.config.Files)
		if len(content) > 0 {
			if err := ioutil.WriteFile(fmt.Sprintf("%s/%s_%s/DEBIAN/conffiles", b.workingdirectory, b.config.Name, b.config.Version), content, 0644); err != nil {
				return err
			}
			logger.Debugf("debian conffiles written")
		} else {
			logger.Debugf("no found files, skip writing")
		}
		return nil
	}, func() error {
		return nil
	})
}

func conffilesContent(files []debian_config.File) []byte {
	buffer := bytes.NewBufferString("")
	for _, file := range files {
		if strings.Index(file.Target, "/etc") == 0 {
			fmt.Fprintf(buffer, "%s\n", file.Target)
		}
	}
	return buffer.Bytes()
}

func (b *builder) createDebianPackageCommand() command.Command {
	return command_adapter.New(func() error {
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

func (b *builder) copyDebianPackageCommand() command.Command {
	return command_adapter.New(func() error {
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

func (b *builder) copyFilesToWorkingDirectoryCommand() command.Command {
	return command_adapter.New(func() error {
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

func (b *builder) cleanWorkingDirectoryCommand() command.Command {
	return command_adapter.New(func() error {
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
		return filename[:pos+1], nil
	}
	return "", fmt.Errorf("can't determine directory of file %s", filename)
}
