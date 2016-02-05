package main

import (
	"flag"
	"io"
	"os"
	"runtime"

	"fmt"

	"github.com/bborbe/debian_utils/zip_extractor"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_ZIP      = "zip"
	PARAMETER_TARGET   = "target"
)

type ExtractZipFile func(filename string, targetDir string) error

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	zipPtr := flag.String(PARAMETER_ZIP, "", "zip")
	targetPtr := flag.String(PARAMETER_TARGET, "", "target")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	z := zip_extractor.New()

	writer := os.Stdout
	err := do(writer, z.ExtractZipFile, *zipPtr, *targetPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, extractZipFile ExtractZipFile, zip string, target string) error {
	if len(zip) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_ZIP)
	}
	if len(target) == 0 {
		return fmt.Errorf("parameter %s missing", PARAMETER_TARGET)
	}
	return extractZipFile(zip, target)
}
