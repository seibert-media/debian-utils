package main

import (
	"flag"
	"io"
	"os"
	"runtime"

	debian_command_list "github.com/bborbe/command/list"
	debian_config "github.com/bborbe/debian_utils/config"
	debian_config_builder "github.com/bborbe/debian_utils/config_builder"
	debian_config_parser "github.com/bborbe/debian_utils/config_parser"
	debian_copier "github.com/bborbe/debian_utils/copier"
	debian_package_creator "github.com/bborbe/debian_utils/package_creator"
	debian_tar_gz_extractor "github.com/bborbe/debian_utils/tar_gz_extractor"
	debian_zip_extractor "github.com/bborbe/debian_utils/zip_extractor"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/bborbe/log"
)

var logger = log.DefaultLogger

const (
	PARAMETER_NAME     = "name"
	PARAMETER_VERSION  = "version"
	PARAMETER_SOURCE   = "source"
	PARAMETER_TARGET   = "target"
	PARAMETER_LOGLEVEL = "loglevel"
	PARAMETER_CONFIG   = "config"
)

type ConfigBuilderWithConfig func(config *debian_config.Config) debian_config_builder.ConfigBuilder

func main() {
	defer logger.Close()
	logLevelPtr := flag.String(PARAMETER_LOGLEVEL, log.INFO_STRING, log.FLAG_USAGE)
	configPtr := flag.String(PARAMETER_CONFIG, "", "config")
	namePtr := flag.String(PARAMETER_NAME, "", "name")
	versionPtr := flag.String(PARAMETER_VERSION, "", "version")
	sourcePtr := flag.String(PARAMETER_SOURCE, "", "source")
	targetPtr := flag.String(PARAMETER_TARGET, "", "target")
	flag.Parse()
	logger.SetLevelThreshold(log.LogStringToLevel(*logLevelPtr))
	logger.Debugf("set log level to %s", *logLevelPtr)

	runtime.GOMAXPROCS(runtime.NumCPU())

	commandListProvider := func() debian_command_list.CommandList {
		return debian_command_list.New()
	}
	configBuilderWithConfig := func(config *debian_config.Config) debian_config_builder.ConfigBuilder {
		return debian_config_builder.NewWithConfig(config)
	}
	config_parser := debian_config_parser.New()
	copier := debian_copier.New()
	zipExtractor := debian_zip_extractor.New()
	tarGzExtractor := debian_tar_gz_extractor.New()
	httpClientBuilder := http_client_builder.New().WithoutProxy()
	httpClient := httpClientBuilder.Build()
	requestbuilderProvider := http_requestbuilder.NewHttpRequestBuilderProvider()
	debianPackageCreator := debian_package_creator.New(commandListProvider, copier, tarGzExtractor.ExtractTarGz, zipExtractor.ExtractZip, httpClient.Do, requestbuilderProvider.NewHttpRequestBuilder)

	writer := os.Stdout
	err := do(writer, config_parser, configBuilderWithConfig, debianPackageCreator, *configPtr, *namePtr, *versionPtr, *sourcePtr, *targetPtr)
	if err != nil {
		logger.Fatal(err)
		logger.Close()
		os.Exit(1)
	}
}

func do(writer io.Writer,
	config_parser debian_config_parser.ConfigParser,
	configBuilderWithConfig ConfigBuilderWithConfig,
	package_creator debian_package_creator.PackageCreator,
	configpath string,
	name string,
	version string,
	source string,
	target string) error {
	var err error
	config := debian_config.DefaultConfig()
	if len(configpath) > 0 {
		if config, err = config_parser.ParseFileToConfig(config, configpath); err != nil {
			return err
		}
	}
	config_builder := configBuilderWithConfig(config)
	config_builder.AddFile(source, target)
	config_builder.Name(name)
	config_builder.Version(version)
	return package_creator.CreatePackage(config_builder.Build())
}
