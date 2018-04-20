package latest_package_creator

import (
	"io"
	"net/http"

	debian_config "github.com/seibert-media/debian-utils/config"
)

type LatestConfluenceTarGzURL func() (string, error)

type LatestConfluenceVersion func() (string, error)

type CreatePackageByReader func(tarGzReader io.Reader, config *debian_config.Config, sourceDir string, targetDir string) error

type Download func(url string) (resp *http.Response, err error)

type LatestDebianPackageCreator interface {
	CreateLatestDebianPackage(config *debian_config.Config, sourceDir string, targetDir string) error
}

type latestDebianPackageCreator struct {
	latestConfluenceTarGzURL LatestConfluenceTarGzURL
	latestConfluenceVersion  LatestConfluenceVersion
	createPackageByReader    CreatePackageByReader
	download                 Download
}

func New(download Download, latestConfluenceTarGzURL LatestConfluenceTarGzURL, latestConfluenceVersion LatestConfluenceVersion, createPackageByReader CreatePackageByReader) *latestDebianPackageCreator {
	l := new(latestDebianPackageCreator)
	l.latestConfluenceTarGzURL = latestConfluenceTarGzURL
	l.latestConfluenceVersion = latestConfluenceVersion
	l.createPackageByReader = createPackageByReader
	l.download = download
	return l
}

func (l *latestDebianPackageCreator) CreateLatestDebianPackage(config *debian_config.Config, sourceDir string, targetDir string) error {
	url, err := l.latestConfluenceTarGzURL()
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
