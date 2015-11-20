package main

import (
	"flag"
	"io"
	"os"
	"runtime"
	"github.com/bborbe/log"
	debian_command_list "github.com/bborbe/debian/command_list"
	debian_config_builder "github.com/bborbe/debian/config_builder"
	debian_package_creator "github.com/bborbe/debian/package_creator"
	debian_copier "github.com/bborbe/debian/copier"
)

var logger = log.DefaultLogger

const (
	PARAMETER_NAME = "name"
	PARAMETER_VERSION = "version"
	PARAMETER_SOURCE = "source"
	PARAMETER_TARGET = "target"
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

	copier := debian_copier.New()
	config_builder := debian_config_builder.New()
	package_creator := debian_package_creator.New(func() debian_command_list.CommandList {
		return debian_command_list.New()
	}, copier)

	writer := os.Stdout
	err := do(writer, config_builder, package_creator, *namePtr, *versionPtr, *sourcePtr, *targetPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, config_builder debian_config_builder.ConfigBuilder, package_creator debian_package_creator.PackageCreator, name string, version string, source string, target string) error {
	logger.Debugf("create deb %s_%s.deb", name, version)
	var err error
	if err = config_builder.AddFile(source, target); err != nil {
		return err
	}
	if err = config_builder.Name(name); err != nil {
		return err
	}
	if err = config_builder.Version(version); err != nil {
		return err
	}
	return package_creator.CreatePackage(config_builder.Build())
}
