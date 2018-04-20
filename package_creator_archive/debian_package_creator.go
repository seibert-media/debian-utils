package package_creator_archive

import (
	"io"
	"os"

	debian_config "github.com/seibert-media/debian-utils/config"
	"github.com/golang/glog"
)

type DebianPackageCreator interface {
	CreatePackage(archivePath string, config *debian_config.Config, sourceDir string, targetDir string) error
}

type CreatePackage func(tarGzReader io.Reader, config *debian_config.Config, sourceDir string, targetDir string) error

type debianPackageCreator struct {
	createPackage CreatePackage
}

func New(createPackage CreatePackage) *debianPackageCreator {
	d := new(debianPackageCreator)
	d.createPackage = createPackage
	return d
}

func (d *debianPackageCreator) CreatePackage(archivePath string, config *debian_config.Config, sourceDir string, targetDir string) error {
	glog.V(2).Infof("CreatePackage with archive %s and version: %s", archivePath, config.Version)
	f, err := os.OpenFile(archivePath, os.O_RDONLY, 0444)
	defer f.Close()
	if err != nil {
		return err
	}
	return d.createPackage(f, config, sourceDir, targetDir)
}
