package main

import (
	"flag"
	"os"
	"runtime"

	"fmt"

	"github.com/bborbe/debian_utils/zip_extractor"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	parameterLoglevel = "loglevel"
	parameterZip      = "zip"
	parameterTarget   = "target"
)

type ExtractZipFile func(filename string, targetDir string) error

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(parameterLoglevel, log.INFO_STRING, log.FLAG_USAGE)
	zipPtr := flag.String(parameterZip, "", "zip")
	targetPtr := flag.String(parameterTarget, "", "target")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	z := zip_extractor.New()

	err := do(z.ExtractZipFile, *zipPtr, *targetPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(extractZipFile ExtractZipFile, zip string, target string) error {
	if len(zip) == 0 {
		return fmt.Errorf("parameter %s missing", parameterZip)
	}
	if len(target) == 0 {
		return fmt.Errorf("parameter %s missing", parameterTarget)
	}
	return extractZipFile(zip, target)
}
