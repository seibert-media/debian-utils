package apt_source_has_changed

import (
	"bufio"
	"io"
	"os"

	"github.com/golang/glog"
)

type AptSourceHasChanged interface {
	HasFileChanged(path string) (bool, error)
}

type aptSourceHasChanged struct {
	hasLineChanged HasLineChanged
}

type HasLineChanged func(line string) (bool, error)

type ReadString func(delim byte) (line string, err error)

func New(hasLineChanged HasLineChanged) *aptSourceHasChanged {
	a := new(aptSourceHasChanged)
	a.hasLineChanged = hasLineChanged
	return a
}

func (a *aptSourceHasChanged) HasFileChanged(path string) (bool, error) {
	glog.V(2).Infof("hasFileChanged - path: %s", path)
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	defer file.Close()
	return hasFileChanged(bufio.NewReader(file).ReadString, a.hasLineChanged)
}

func hasFileChanged(readString ReadString, hasLineChanged HasLineChanged) (bool, error) {
	for {
		var err error
		var line string
		var changed bool
		if line, err = readString('\n'); err != nil {
			if err == io.EOF {
				if changed, err = hasLineChanged(line); err != nil {
					return false, err
				}
				if changed {
					glog.V(2).Infof("line has changed => true")
					return true, nil
				}
			}
			return false, err
		}
		if changed, err = hasLineChanged(line); err != nil {
			return false, err
		}
		if changed {
			glog.V(2).Infof("line has changed => true")
			return true, nil
		}
	}
}
