package main

import (
	"flag"
	"io"
	"os"
	"runtime"

	debian_apt_source_has_changed "github.com/bborbe/debian_utils/apt_source_has_changed"
	debian_line_inspector "github.com/bborbe/debian_utils/apt_source_line_inspector"
	debian_url_downloader "github.com/bborbe/debian_utils/url_downloader"

	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_PATH     = "path"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	pathPtr := flag.String(PARAMETER_PATH, "", "path")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	httpClientBuilder := http_client_builder.New().WithoutProxy()
	httpClient := httpClientBuilder.Build()
	requestbuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	downloader := debian_url_downloader.New(httpClient.Do, requestbuilderProvider.NewHttpRequestBuilder)
	lineInspector := debian_line_inspector.New(downloader.DownloadUrl)
	hasChanged := debian_apt_source_has_changed.New(lineInspector.HasLineChanged)

	writer := os.Stdout
	bool, err := do(writer, hasChanged, *pathPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
	if bool {
		logger.Close()
		os.Exit(0)
	} else {
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, hasChanged debian_apt_source_has_changed.AptSourceHasChanged, path string) (bool, error) {
	logger.Debugf("update repos in apt source list: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, err
	}
	return hasChanged.HasFileChanged(path)
}
