package main

import (
	"flag"
	"io"
	"os"
	"runtime"

	"github.com/bborbe/debian_utils/apt_source_list_updater"
	"github.com/bborbe/debian_utils/url_downloader"
	http_client "github.com/bborbe/http/client"
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
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	pathPtr := flag.String(PARAMETER_PATH, "", "path")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	client := http_client.GetClientWithoutProxy()
	requestbuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	downloader := url_downloader.New(client, requestbuilderProvider.NewHttpRequestBuilder)
	updater := apt_source_list_updater.New(downloader.DownloadUrl)

	writer := os.Stdout
	err := do(writer, updater, *pathPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, updater apt_source_list_updater.AptSourceListUpdater, path string) error {
	logger.Debugf("update repos in apt source list: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	return updater.UpdateAptSourceList(path)
}
