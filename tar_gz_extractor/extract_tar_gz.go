package tar_gz_extractor

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bborbe/log"
)

type TarGzExtractor interface {
	ExtractTarGz(fileReader io.Reader, targetDir string) error
}

type tarGzExtractor struct {
}

func New() *tarGzExtractor {
	return new(tarGzExtractor)
}

var logger = log.DefaultLogger

func (e *tarGzExtractor) ExtractTarGz(fileReader io.Reader, targetDir string) error {
	logger.Debugf("extract tar fz to %s", targetDir)

	gw, err := gzip.NewReader(fileReader)
	if err != nil {
		return err
	}
	defer gw.Close()

	tr := tar.NewReader(gw)
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		path := fmt.Sprintf("%s/%s", targetDir, hdr.Name)
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
