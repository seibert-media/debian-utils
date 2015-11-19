package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"runtime"

	debian_command_list "github.com/bborbe/debian/command_list"
	debian_config "github.com/bborbe/debian/config"
	debian_config_parser "github.com/bborbe/debian/config_parser"
	debian_package_creator "github.com/bborbe/debian/package_creator"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	PARAMETER_CONFIG   = "config"
	PARAMETER_LOGLEVEL = "loglevel"
)

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, "one of OFF,TRACE,DEBUG,INFO,WARN,ERROR")
	configPtr := flag.String(PARAMETER_CONFIG, "", "name")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	config_parser := debian_config_parser.New()
	package_creator := debian_package_creator.New(func() debian_command_list.CommandList {
		return debian_command_list.New()
	})

	writer := os.Stdout
	err := do(writer, config_parser, package_creator, *configPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer, config_parser debian_config_parser.ConfigParser, package_creator debian_package_creator.PackageCreator, configpath string) error {
	logger.Debugf("create deb by config %s", configpath)

	if len(configpath) == 0 {
		return fmt.Errorf("parameter config missing")
	}
	var err error
	var content []byte
	var config *debian_config.Config
	if content, err = ioutil.ReadFile(configpath); err != nil {
		return err
	}
	if config, err = config_parser.ParseConfig(content); err != nil {
		return err
	}
	return package_creator.CreatePackage(config)
}
