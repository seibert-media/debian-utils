package zip_extractor

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/bborbe/log"
)

type ZipExtractor interface {
	ExtractZip(fileReader io.Reader, targetDir string) error
}

type zipExtractor struct {
}

func New() *zipExtractor {
	return new(zipExtractor)
}

var logger = log.DefaultLogger

func (e *zipExtractor) ExtractZip(fileReader io.Reader, targetDir string) error {
	logger.Debugf("extract zip")

	filename := "/tmp/test.zip"
	defer os.Remove(filename)
	err := write(fileReader, filename)
	if err != nil {
		return err
	}
	return e.ExtractZipFile(filename, targetDir)
}

func (e *zipExtractor) ExtractZipFile(filename string, targetDir string) error {
	logger.Debugf("extract zip %s", filename)
	z, err := zip.OpenReader(filename)
	if err != nil {
		return err
	}
	for _, f := range z.File {
		path := fmt.Sprintf("%s/%s", targetDir, f.Name)
		if f.FileInfo().IsDir() {
			logger.Debugf("extract dir %s", f.Name)
			mkdir(path, f.FileInfo().Mode())
		} else {
			logger.Debugf("extract file %s", f.Name)
			reader, err := f.Open()
			if err != nil {
				return err
			}
			err = extractFile(path, f.FileInfo().Mode(), reader)
			if err != nil {
				return err
			}
			reader.Close()
		}
	}
	logger.Debugf("zip extracted %s", filename)
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
	ow, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_TRUNC, mode)
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

func write(fileReader io.Reader, filename string) error {
	out, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, fileReader)
	if err != nil {
		return err
	}
	return nil
}
