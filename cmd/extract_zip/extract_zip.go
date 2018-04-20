package main

import (
	"flag"
	"runtime"

	"fmt"

	"github.com/seibert-media/debian-utils/zip_extractor"
	"github.com/golang/glog"
)

const (
	parameterZip    = "zip"
	parameterTarget = "target"
)

type ExtractZipFile func(filename string, targetDir string) error

var (
	zipPtr    = flag.String(parameterZip, "", "zip")
	targetPtr = flag.String(parameterTarget, "", "target")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
	runtime.GOMAXPROCS(runtime.NumCPU())

	z := zip_extractor.New()

	err := do(
		z.ExtractZipFile,
		*zipPtr,
		*targetPtr,
	)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	extractZipFile ExtractZipFile,
	zip string,
	target string,
) error {
	glog.Infof("zip: %v target: %v", zip, target)
	if len(zip) == 0 {
		return fmt.Errorf("parameter %s missing", parameterZip)
	}
	if len(target) == 0 {
		return fmt.Errorf("parameter %s missing", parameterTarget)
	}
	return extractZipFile(zip, target)
}
