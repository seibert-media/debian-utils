package main

import (
	"flag"
	"io"
	"os"
	"runtime"

	debian_apt_source_list_updater "github.com/bborbe/debian_utils/apt_source_list_updater"
	debian_url_downloader "github.com/bborbe/debian_utils/url_downloader"
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
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	pathPtr := flag.String(PARAMETER_PATH, "", "path")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	client := http_client.GetClientWithoutProxy()
	requestbuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	downloader := debian_url_downloader.New(client, requestbuilderProvider.NewHttpRequestBuilder)
	updater := debian_apt_source_list_updater.New(downloader.DownloadUrl)

	writer := os.Stdout
	err := do(writer, updater, *pathPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, updater debian_apt_source_list_updater.AptSourceListUpdater, path string) error {
	logger.Debugf("update repos in apt source list: %s", path)
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return err
	}
	return updater.UpdateAptSourceList(path)
}
