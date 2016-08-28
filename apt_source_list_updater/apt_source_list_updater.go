package apt_source_list_updater

import (
	"os"

	"fmt"
	"os/exec"

	"path/filepath"

	"github.com/golang/glog"
)

type AptSourceListUpdater interface {
	UpdateAptSourceList(path string) error
}

type aptSourceListUpdater struct {
	hasFileChanged HasFileChanged
}

type HasFileChanged func(path string) (bool, error)

func New(hasFileChanged HasFileChanged) *aptSourceListUpdater {
	a := new(aptSourceListUpdater)
	a.hasFileChanged = hasFileChanged
	return a
}

func (a *aptSourceListUpdater) UpdateAptSourceList(path string) error {
	glog.V(2).Infof("UpdateAptSourceList - path: %s", path)
	changed, err := a.hasFileChanged(path)
	if err != nil {
		return err
	}
	if changed {
		glog.V(2).Infof("has changed => trigger update")
		return a.update(path)
	}
	glog.V(2).Infof("nothing has changed")
	return nil
}

func (a *aptSourceListUpdater) update(path string) error {
	path, err := filepath.Abs(path)
	if err != nil {
		return err
	}
	args := []string{"update", "-o", fmt.Sprintf("Dir::Etc::sourcelist=%s", path), "-o", "Dir::Etc::sourceparts=-", "-o", "APT::Get::List-Cleanup=0"}
	cmd := exec.Command("apt-get", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
