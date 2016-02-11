package apt_source_has_changed

import (
	"bufio"
	"io"
	"os"

	"fmt"
	"regexp"
	"runtime"
	"strings"

	"io/ioutil"

	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

type AptSourceHasChanged interface {
	HasFileChanged(path string) (bool, error)
}

type aptSourceHasChanged struct {
	downloadUrl DownloadUrl
}

type DownloadUrl func(url string) (string, error)

func New(downloadUrl DownloadUrl) *aptSourceHasChanged {
	a := new(aptSourceHasChanged)
	a.downloadUrl = downloadUrl
	return a
}

func (a *aptSourceHasChanged) HasFileChanged(path string) (bool, error) {
	logger.Debugf("hasFileChanged - path: %s", path)
	file, err := os.Open(path)
	if err != nil {
		return false, err
	}
	in := bufio.NewReader(file)
	for {
		var line string
		var changed bool
		if line, err = in.ReadString('\n'); err != nil {
			if err == io.EOF {
				return false, nil
			}
			return false, err
		}
		if changed, err = a.hasLineChanged(line); err != nil {
			return false, err
		}
		if changed {
			logger.Debugf("line has changed")
			return true, nil
		}
	}
	return false, nil
}

func (a *aptSourceHasChanged) hasLineChanged(line string) (bool, error) {
	logger.Debugf("updateLine - line: %s", line)
	if strings.Index(line, "deb ") != 0 {
		return false, nil
	}
	infos, err := ParseLine(line)
	if err != nil {
		return false, err
	}
	return a.hasPackages(infos)
}

type infos struct {
	url          string
	distribution string
	architecture string
	component    string
}

func ParseLine(line string) (*infos, error) {
	i := new(infos)
	{
		re := regexp.MustCompile(`deb\s+\[arch=(.*?)\]\s+([^\s]+)\s+([^\s]+)\s+([^\s]+)`)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			i.architecture = matches[1]
			i.url = matches[2]
			i.distribution = matches[3]
			i.component = matches[4]
			return i, nil
		}
	}
	{
		re := regexp.MustCompile(`deb\s+([^\s]+)\s([^\s]+)\s+([^\s]+)`)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			i.architecture = runtime.GOARCH
			i.url = matches[1]
			i.distribution = matches[2]
			i.component = matches[3]
			return i, nil
		}
	}
	{
		re := regexp.MustCompile(`deb\s+([^\s]+)\s([^\s]+)`)
		matches := re.FindStringSubmatch(line)
		if len(matches) > 1 {
			i.architecture = runtime.GOARCH
			i.url = matches[1]
			i.distribution = matches[2]
			return i, nil
		}
	}
	return nil, fmt.Errorf("parse line failed: %s", line)
}

func (a *aptSourceHasChanged) hasPackages(infos *infos) (bool, error) {
	remotePackagesUrl := infos.RemotePackagesUrl()
	logger.Debugf("remote packages url: %s", remotePackagesUrl)
	remotePackagesContent, err := a.downloadUrl(remotePackagesUrl)
	if err != nil {
		return false, err
	}
	localPackagesFile := infos.LocalPackagesFile()
	logger.Debugf("local packages file: %s", localPackagesFile)
	localPackagesContent, err := a.readFile(localPackagesFile)
	if err != nil {
		// return false, err
		// return true if file not found
		return true, nil
	}
	return localPackagesContent != remotePackagesContent, nil
}

func (i *infos) RemotePackagesUrl() string {
	return fmt.Sprintf("%s/dists/%s/%s/binary-%s/Packages", i.url, i.distribution, i.component, i.architecture)
}

func (i *infos) LocalPackagesFile() string {
	atPos := strings.Index(i.url, "@")
	var host string
	if atPos != -1 {
		host = i.url[atPos+1:]
	} else {
		pos := strings.Index(i.url, "://")
		host = i.url[pos+3:]
	}
	return fmt.Sprintf("/var/lib/apt/lists/%s_dists_%s_%s_binary-%s_Packages", strings.Replace(host, "/", "_", -1), i.distribution, i.component, i.architecture)
}

func (a *aptSourceHasChanged) readFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
