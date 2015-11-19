package main

import (
	"flag"
	"os"

	"runtime"

	"io"

	"github.com/bborbe/debian/command_list"
	"github.com/bborbe/debian/package_builder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	PARAMETER_NAME     = "name"
	PARAMETER_VERSION  = "version"
	PARAMETER_SOURCE   = "source"
	PARAMETER_TARGET   = "target"
	PARAMETER_LOGLEVEL = "loglevel"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	namePtr := flag.String(PARAMETER_NAME, "", "name")
	versionPtr := flag.String(PARAMETER_VERSION, "", "version")
	sourcePtr := flag.String(PARAMETER_SOURCE, "", "source")
	targetPtr := flag.String(PARAMETER_TARGET, "", "target")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	builder := package_creator.New(command_list.New())

	writer := os.Stdout
	err := do(writer, builder, *namePtr, *versionPtr, *sourcePtr, *targetPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, builder package_creator.Builder, name string, version string, source string, target string) error {
	logger.Debugf("create deb %s_%s.deb", name, version)
	if err := builder.AddFile(source, target); err != nil {
		return err
	}
	if err := builder.Name(name); err != nil {
		return err
	}
	if err := builder.Version(version); err != nil {
		return err
	}
	return builder.Build()
}
