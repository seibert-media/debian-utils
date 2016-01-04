package latest_package_creator

import (
	"io"
	"net/http"

	debian_config "github.com/bborbe/debian_utils/config"
)

type LatestConfluenceTarGzUrl func() (string, error)

type LatestConfluenceVersion func() (string, error)

type CreatePackageByReader func(tarGzReader io.Reader, config *debian_config.Config, sourceDir string, targetDir string) error

type Download func(url string) (resp *http.Response, err error)

type LatestDebianPackageCreator interface {
	CreateLatestConfluenceDebianPackage(config *debian_config.Config, sourceDir string, targetDir string) error
}

type latestDebianPackageCreator struct {
	latestConfluenceTarGzUrl LatestConfluenceTarGzUrl
	latestConfluenceVersion  LatestConfluenceVersion
	createPackageByReader    CreatePackageByReader
	download                 Download
}

func New(download Download, latestConfluenceTarGzUrl LatestConfluenceTarGzUrl, latestConfluenceVersion LatestConfluenceVersion, createPackageByReader CreatePackageByReader) *latestDebianPackageCreator {
	l := new(latestDebianPackageCreator)
	l.latestConfluenceTarGzUrl = latestConfluenceTarGzUrl
	l.latestConfluenceVersion = latestConfluenceVersion
	l.createPackageByReader = createPackageByReader
	l.download = download
	return l
}

func (l *latestDebianPackageCreator) CreateLatestConfluenceDebianPackage(config *debian_config.Config, sourceDir string, targetDir string) error {
	url, err := l.latestConfluenceTarGzUrl()
	if err != nil {
		return err
	}
	config.Version, err = l.latestConfluenceVersion()
	if err != nil {
		return err
	}
	resp, err := l.download(url)
	if err != nil {
		return err
	}
	return l.createPackageByReader(resp.Body, config, sourceDir, targetDir)
}
