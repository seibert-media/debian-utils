package copier

import (
	"fmt"
	"io"
	"os"

	"github.com/golang/glog"
)

type Copier interface {
	Copy(source string, destination string) error
}

type copier struct {
}

func New() *copier {
	return new(copier)
}

func (c *copier) Copy(source string, target string) error {
	glog.V(2).Infof("Copy %s => %s", source, target)
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
	glog.V(2).Infof("CopyDir %s => %s", source, target)
	finfo, err := os.Stat(source)
	if err != nil {
		glog.V(2).Infof("get stat failed")
		return err
	}
	if err = os.MkdirAll(target, finfo.Mode()); err != nil {
		glog.V(2).Infof("mkdir target dir failed")
		return err
	}
	dir, err := os.Open(source)
	defer dir.Close()
	if err != nil {
		glog.V(2).Infof("open source failed")
		return err
	}
	files, err := dir.Readdir(-1)
	if err != nil {
		glog.V(2).Infof("read source dir failed")
		return err
	}
	for _, file := range files {
		glog.V(2).Infof("file: %s", file.Name())
		if err = c.Copy(fmt.Sprintf("%s/%s", source, file.Name()), fmt.Sprintf("%s/%s", target, file.Name())); err != nil {
			return err
		}
	}
	return nil
}

func (c *copier) CopyFile(source string, target string) error {
	glog.V(2).Infof("CopyFile %s => %s", source, target)
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
