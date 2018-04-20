package main

import (
	"flag"
	"os"
	"runtime"

	debian_apt_source_has_changed "github.com/seibert-media/debian-utils/apt_source_has_changed"
	debian_line_inspector "github.com/seibert-media/debian-utils/apt_source_line_inspector"
	debian_apt_source_list_updater "github.com/seibert-media/debian-utils/apt_source_list_updater"
	debian_url_downloader "github.com/seibert-media/debian-utils/url_downloader"

	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/golang/glog"
)

const (
	parameterPath = "path"
)

var (
	pathPtr = flag.String(parameterPath, "", "path")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClientBuilder := http_client_builder.New().WithoutProxy()
	httpClient := httpClientBuilder.Build()
	requestbuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	downloader := debian_url_downloader.New(httpClient.Do, requestbuilderProvider.NewHTTPRequestBuilder)
	lineInspector := debian_line_inspector.New(downloader.DownloadURL)
	hasChanged := debian_apt_source_has_changed.New(lineInspector.HasLineChanged)
	updater := debian_apt_source_list_updater.New(hasChanged.HasFileChanged)

	err := do(
		updater,
		*pathPtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	updater debian_apt_source_list_updater.AptSourceListUpdater,
	path string,
) error {
	glog.Infof("update repos in apt source list: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	return updater.UpdateAptSourceList(path)
}
