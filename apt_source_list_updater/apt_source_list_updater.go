package apt_source_list_updater

import (
	"os"

	"fmt"
	"os/exec"

	"github.com/bborbe/log"
)

type AptSourceListUpdater interface {
	UpdateAptSourceList(path string) error
}

type aptSourceListUpdater struct {
	hasFileChanged HasFileChanged
}

var logger = log.DefaultLogger

type HasFileChanged func(path string) (bool, error)

func New(hasFileChanged HasFileChanged) *aptSourceListUpdater {
	a := new(aptSourceListUpdater)
	a.hasFileChanged = hasFileChanged
	return a
}

func (a *aptSourceListUpdater) UpdateAptSourceList(path string) error {
	logger.Debugf("UpdateAptSourceList - path: %s", path)
	changed, err := a.hasFileChanged(path)
	if err != nil {
		return err
	}
	if changed {
		logger.Debugf("has changed => trigger update")
		return a.update(path)
	}
	logger.Debugf("nothing has changed")
	return nil
}

func (a *aptSourceListUpdater) update(path string) error {
	args := []string{"update", "-o", fmt.Sprintf("Dir::Etc::sourcelist=%s", path), "-o", "Dir::Etc::sourceparts=-", "-o", "APT::Get::List-Cleanup=0"}
	cmd := exec.Command("apt-get", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
