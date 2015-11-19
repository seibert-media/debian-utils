package package_creator

import (
	"io"
	"io/ioutil"
	"os"

	"fmt"
	"strings"

	"bytes"
	"os/exec"

	"github.com/bborbe/debian/command"
	"github.com/bborbe/debian/command_list"
	"github.com/bborbe/log"
)

type Builder interface {
	Build() error
	Name(name string) error
	Version(version string) error
	Section(section string) error
	Priority(priority string) error
	Architecture(architecture string) error
	Maintainer(maintainer string) error
	Description(description string) error
	AddFile(source string, target string) error
}

type builder struct {
	name             string
	version          string
	description      string
	section          string
	priority         string
	architecture     string
	maintainer       string
	files            map[string]string
	commandList      command_list.CommandList
	workingdirectory string
}

var logger = log.DefaultLogger

func New(commandList command_list.CommandList) *builder {
	b := new(builder)
	b.commandList = commandList
	b.section = "base"
	b.priority = "optional"
	b.architecture = "all"
	b.maintainer = "Benjamin Borbe <bborbe@rocketnews.de>"
	b.description = "-"
	b.files = map[string]string{}
	return b
}

func (b *builder) Section(section string) error {
	logger.Debugf("Section %s", section)
	if len(section) == 0 {
		return fmt.Errorf("section empty")
	}
	b.section = section
	return nil
}

func (b *builder) Priority(priority string) error {
	logger.Debugf("Priority %s", priority)
	if len(priority) == 0 {
		return fmt.Errorf("priority empty")
	}
	b.priority = priority
	return nil
}

func (b *builder) Architecture(architecture string) error {
	logger.Debugf("Architecture %s", architecture)
	if len(architecture) == 0 {
		return fmt.Errorf("architecture empty")
	}
	b.architecture = architecture
	return nil
}

func (b *builder) Maintainer(maintainer string) error {
	logger.Debugf("Maintainer %s", maintainer)
	if len(maintainer) == 0 {
		return fmt.Errorf("maintainer empty")
	}
	b.maintainer = maintainer
	return nil
}

func (b *builder) Description(description string) error {
	logger.Debugf("Description %s", description)
	if len(description) == 0 {
		return fmt.Errorf("description empty")
	}
	b.description = description
	return nil
}

func (b *builder) Name(name string) error {
	logger.Debugf("Name %s", name)
	if len(name) == 0 {
		return fmt.Errorf("name empty")
	}
	b.name = name
	return nil
}

func (b *builder) Version(version string) error {
	logger.Debugf("Version %s", version)
	if len(version) == 0 {
		return fmt.Errorf("version empty")
	}
	b.version = version
	return nil
}

func (b *builder) AddFile(source string, target string) error {
	if len(source) == 0 {
		return fmt.Errorf("source empty")
	}
	if len(target) == 0 {
		return fmt.Errorf("target empty")
	}
	b.files[target] = source
	return nil
}

func (b *builder) Build() error {
	logger.Debug("Build")
	b.commandList.Add(b.validateCommand())
	b.commandList.Add(b.createWorkingDirectoryCommand())
	b.commandList.Add(b.createDebianFolderCommand())
	b.commandList.Add(b.createDebianControlCommand())
	b.commandList.Add(b.copyFilesToWorkingDirectoryCommand())
	b.commandList.Add(b.createDebianPackageCommand())
	b.commandList.Add(b.copyDebianPackageCommand())
	b.commandList.Add(b.cleanWorkingDirectoryCommand())
	return b.commandList.Run()
}

func (b *builder) validateCommand() command.Command {
	return command.New(func() error {
		logger.Debug("validate")
		if len(b.files) == 0 {
			return fmt.Errorf("add at least one file")
		}
		if len(b.name) == 0 {
			return fmt.Errorf("name missing")
		}
		if len(b.version) == 0 {
			return fmt.Errorf("version missing")
		}
		logger.Debug("validate success")
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) createWorkingDirectoryCommand() command.Command {
	return command.New(func() error {
		logger.Debugf("create working directory")
		var err error
		if b.workingdirectory, err = ioutil.TempDir("", fmt.Sprintf("%s_%s", b.name, b.version)); err != nil {
			return err
		}
		logger.Debugf("working directory %s create", b.workingdirectory)
		return nil
	}, func() error {
		return os.RemoveAll(b.workingdirectory)
	})
}

func (b *builder) createDebianFolderCommand() command.Command {
	return command.New(func() error {
		path := fmt.Sprintf("%s/%s_%s/DEBIAN", b.workingdirectory, b.name, b.version)
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
	return command.New(func() error {
		logger.Debugf("write debian control")
		if err := ioutil.WriteFile(fmt.Sprintf("%s/%s_%s/DEBIAN/control", b.workingdirectory, b.name, b.version), b.controlContent(), 0644); err != nil {
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
	fmt.Fprintf(buffer, "Package: %s\n", b.name)
	fmt.Fprintf(buffer, "Version: %s\n", b.version)
	fmt.Fprintf(buffer, "Section: %s\n", b.section)
	fmt.Fprintf(buffer, "Priority: %s\n", b.priority)
	fmt.Fprintf(buffer, "Architecture: %s\n", b.architecture)
	fmt.Fprintf(buffer, "Maintainer: %s\n", b.maintainer)
	fmt.Fprintf(buffer, "Description: %s\n", b.description)
	return buffer.Bytes()
}

func (b *builder) createDebianPackageCommand() command.Command {
	return command.New(func() error {
		logger.Debugf("create debian package")
		cmd := exec.Command("dpkg-deb", "--build", fmt.Sprintf("%s_%s", b.name, b.version))
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
	return command.New(func() error {
		logger.Debugf("copy debian package")
		var dir string
		var err error
		if dir, err = os.Getwd(); err != nil {
			return err
		}
		if err = copy(fmt.Sprintf("%s/%s_%s.deb", dir, b.name, b.version), fmt.Sprintf("%s/%s_%s.deb", b.workingdirectory, b.name, b.version)); err != nil {
			return err
		}
		logger.Debugf("debian package copied")
		return nil
	}, func() error {
		return nil
	})
}

func (b *builder) copyFilesToWorkingDirectoryCommand() command.Command {
	return command.New(func() error {
		logger.Debugf("copy files")
		for target, source := range b.files {
			var err error
			var directory string
			var filename string
			filename = fmt.Sprintf("%s/%s_%s/%s", b.workingdirectory, b.name, b.version, target)
			if directory, err = dirOf(filename); err != nil {
				return err
			}
			if err = createDirectory(directory); err != nil {
				return err
			}
			if err = copy(filename, source); err != nil {
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
	return command.New(func() error {
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

func copy(dst, src string) error {
	finfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(dst, os.O_RDWR|os.O_CREATE|os.O_TRUNC, finfo.Mode())
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}
