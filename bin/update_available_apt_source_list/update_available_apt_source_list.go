package main

import (
	"flag"
	"os"
	"runtime"

	debian_apt_source_has_changed "github.com/bborbe/debian_utils/apt_source_has_changed"
	debian_line_inspector "github.com/bborbe/debian_utils/apt_source_line_inspector"
	debian_url_downloader "github.com/bborbe/debian_utils/url_downloader"

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

	bool, err := do(hasChanged, *pathPtr)
	if err != nil {
		glog.Exit(err)
	}
	if bool {
		glog.Flush()
		os.Exit(0)
	} else {
		glog.Flush()
		os.Exit(1)
	}
}

func do(hasChanged debian_apt_source_has_changed.AptSourceHasChanged, path string) (bool, error) {
	glog.V(2).Infof("update repos in apt source list: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return hasChanged.HasFileChanged(path)
}
