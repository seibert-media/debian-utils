package line_inspector

import (
	"fmt"
	"regexp"
	"runtime"
	"strings"

	"io/ioutil"

	"github.com/golang/glog"
)

type LineInspector interface {
	HasLineChanged(line string) (bool, error)
}

type lineInspector struct {
	downloadURL DownloadURL
}

type DownloadURL func(url string) (string, error)

func New(downloadURL DownloadURL) *lineInspector {
	a := new(lineInspector)
	a.downloadURL = downloadURL
	return a
}

func (a *lineInspector) HasLineChanged(line string) (bool, error) {
	glog.V(2).Infof("updateLine - line: %s", line)
	if strings.Index(line, "deb ") != 0 {
		return false, nil
	}
	infos, err := parseLine(line)
	if err != nil {
		return false, err
	}
	return a.compareLocalAndRemotePackage(infos)
}

type infos struct {
	url          string
	distribution string
	architecture string
	component    string
}

func parseLine(line string) (*infos, error) {
	glog.V(2).Infof("parse line %s", line)
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

func (a *lineInspector) compareLocalAndRemotePackage(infos *infos) (bool, error) {
	remotePackagesURL := infos.RemotePackagesURL()
	glog.V(2).Infof("remote packages url: %s", remotePackagesURL)
	remotePackagesContent, err := a.downloadURL(remotePackagesURL)
	if err != nil {
		glog.V(2).Infof("fetch remote package failed => return false")
		return false, err
	}
	localPackagesFile := infos.LocalPackagesFile()
	glog.V(2).Infof("local packages file: %s", localPackagesFile)
	localPackagesContent, err := a.readFile(localPackagesFile)
	if err != nil {
		glog.V(2).Infof("read local package failed => return true")
		// return false, err
		// return true if file not found
		return true, nil
	}
	result := localPackagesContent != remotePackagesContent
	glog.V(2).Infof("compare local and remote %v", result)
	return result, nil
}

func (i *infos) RemotePackagesURL() string {
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

func (a *lineInspector) readFile(path string) (string, error) {
	content, err := ioutil.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
