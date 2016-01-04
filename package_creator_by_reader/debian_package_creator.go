package package_creator_by_reader

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/bborbe/command"
	command_adapter "github.com/bborbe/command/adapter"
	command_list "github.com/bborbe/command/list"
	debian_config "github.com/bborbe/debian_utils/config"
	debian_config_builder "github.com/bborbe/debian_utils/config_builder"
	debian_package_creator "github.com/bborbe/debian_utils/package_creator"
	"github.com/bborbe/log"
)

type Creator interface {
	CreatePackage(tarGzReader io.Reader, config *debian_config.Config, sourceDir string, targetDir string) error
}

var logger = log.DefaultLogger

type CommandListProvider func() command_list.CommandList

type creator struct {
	commandListProvider CommandListProvider
	packageCreator      debian_package_creator.PackageCreator
}

func New(commandListProvider CommandListProvider, debianPackageCreator debian_package_creator.PackageCreator) *creator {
	d := new(creator)
	d.commandListProvider = commandListProvider
	d.packageCreator = debianPackageCreator
	return d
}

func (d *creator) CreatePackage(tarGzReader io.Reader, config *debian_config.Config, sourceDir string, targetDir string) error {
	b := new(builder)
	b.packageCreator = d.packageCreator
	b.commandList = d.commandListProvider()
	b.config = config
	b.tarGzReader = tarGzReader
	b.commandList.Add(b.createWorkingDirectoryCommand())
	b.commandList.Add(b.extractTarGzCommand())
	b.commandList.Add(b.createDebianPackageCommand())
	b.commandList.Add(b.removeWorkingDirectoryCommand())
	b.sourceDir = sourceDir
	b.targetDir = targetDir
	return b.commandList.Run()
}

type builder struct {
	commandList      command_list.CommandList
	packageCreator   debian_package_creator.PackageCreator
	workingdirectory string
	tarGzReader      io.Reader
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

func (b *builder) extractTarGzCommand() command.Command {
	return command_adapter.New(func() error {
		logger.Debugf("extract tar fz")

		gw, err := gzip.NewReader(b.tarGzReader)
		defer gw.Close()
		if err != nil {
			return err
		}

		tr := tar.NewReader(gw)
		for {
			hdr, err := tr.Next()
			if err == io.EOF {
				break
			}
			if err != nil {
				return err
			}
			path := fmt.Sprintf("%s/%s", b.workingdirectory, hdr.Name)
			switch hdr.Typeflag {
			case tar.TypeDir:
				if err = mkdir(path, os.FileMode(hdr.Mode)); err != nil {
					return err
				}
			case tar.TypeReg:
				if err = extractFile(path, os.FileMode(hdr.Mode), tr); err != nil {
					return err
				}
			default:
				logger.Debugf("Can't: %c, %s\n", hdr.Typeflag, path)
			}
		}

		logger.Debugf("tar fz extracted")
		return nil
	}, func() error { return nil })
}

func extractFile(path string, mode os.FileMode, tr io.Reader) error {
	logger.Debugf("extract file: %s %v", path, mode)
	dir := filepath.Dir(path)
	_, err := os.Stat(dir)
	if err != nil {
		err := mkdir(dir, os.FileMode(0777))
		if err != nil {
			return err
		}
	}
	ow, err := os.OpenFile(path, os.O_RDWR | os.O_CREATE | os.O_TRUNC, mode)
	defer ow.Close()
	if err != nil {
		logger.Debugf("open file failed: %s %v", path, mode)
		return err
	}
	if err != nil {
		return err
	}
	if _, err := io.Copy(ow, tr); err != nil {
		return err
	}
	return nil
}

func mkdir(path string, mode os.FileMode) error {
	logger.Debugf("mkdir: %s %v", path, mode)
	return os.MkdirAll(path, mode)
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