package copier

import (
	"fmt"
	"io"
	"os"

	"github.com/bborbe/log"
)

type Copier interface {
	Copy(source string, destination string) error
}

type copier struct {
}

var logger = log.DefaultLogger

func New() *copier {
	return new(copier)
}

func (c *copier) Copy(source string, target string) error {
	logger.Debugf("Copy %s => %s", source, target)
	finfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	if finfo.IsDir() {
		return c.CopyDir(source, target)
	}
	return c.CopyFile(source, target)
}

func (c *copier) CopyDir(source string, target string) error {
	logger.Debugf("CopyDir %s => %s", source, target)
	finfo, err := os.Stat(source)
	if err != nil {
		return err
	}
	if err = os.Mkdir(target, finfo.Mode()); err != nil {
		return err
	}
	dir, err := os.Open(source)
	defer dir.Close()
	if err != nil {
		return err
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		return err
	}
	for _, file := range files {
		logger.Debugf("file: %s", file.Name())
		if err = c.Copy(fmt.Sprintf("%s/%s", source, file.Name()), fmt.Sprintf("%s/%s", target, file.Name())); err != nil {
			return err
		}
	}
	return nil
}

func (c *copier) CopyFile(source string, target string) error {
	logger.Debugf("CopyFile %s => %s", source, target)
	finfo, err := os.Stat(source)
	if err != nil {
		return err
	}

	in, err := os.Open(source)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.OpenFile(target, os.O_RDWR|os.O_CREATE|os.O_TRUNC, finfo.Mode())
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
