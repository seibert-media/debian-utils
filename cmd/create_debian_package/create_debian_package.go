package main

import (
	"flag"
	"runtime"

	debian_command_list "github.com/bborbe/command/list"
	debian_config "github.com/seibert-media/debian-utils/config"
	debian_config_builder "github.com/seibert-media/debian-utils/config_builder"
	debian_config_parser "github.com/seibert-media/debian-utils/config_parser"
	debian_copier "github.com/seibert-media/debian-utils/copier"
	debian_package_creator "github.com/seibert-media/debian-utils/package_creator"
	debian_tar_gz_extractor "github.com/seibert-media/debian-utils/tar_gz_extractor"
	debian_zip_extractor "github.com/seibert-media/debian-utils/zip_extractor"
	http_client_builder "github.com/bborbe/http/client_builder"
	http_requestbuilder "github.com/bborbe/http/requestbuilder"
	"github.com/golang/glog"
)

const (
	parameterName    = "name"
	parameterVersion = "version"
	parameterSource  = "source"
	parameterTarget  = "target"
	parameterConfig  = "config"
)

type ConfigBuilderWithConfig func(config *debian_config.Config) debian_config_builder.ConfigBuilder

var (
	configPtr  = flag.String(parameterConfig, "", "config")
	namePtr    = flag.String(parameterName, "", "name")
	versionPtr = flag.String(parameterVersion, "", "version")
	sourcePtr  = flag.String(parameterSource, "", "source")
	targetPtr  = flag.String(parameterTarget, "", "target")
)

func main() {
	defer glog.Flush()
	glog.CopyStandardLogTo("info")
	flag.Parse()
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
	requestbuilderProvider := http_requestbuilder.NewHTTPRequestBuilderProvider()
	debianPackageCreator := debian_package_creator.New(commandListProvider, copier, tarGzExtractor.ExtractTarGz, zipExtractor.ExtractZip, httpClient.Do, requestbuilderProvider.NewHTTPRequestBuilder)

	err := do(config_parser, configBuilderWithConfig, debianPackageCreator, *configPtr, *namePtr, *versionPtr, *sourcePtr, *targetPtr)
	if err != nil {
		glog.Exit(err)
	}
}

func do(
	config_parser debian_config_parser.ConfigParser,
	configBuilderWithConfig ConfigBuilderWithConfig,
	package_creator debian_package_creator.PackageCreator,
	configpath string,
	name string,
	version string,
	source string,
	target string,
) error {
	glog.V(1).Infof("config: %v name: %v version: %v source: %v target: %v", configpath, name, version, source, target)
	var err error
	config := debian_config.DefaultConfig()
	if len(configpath) > 0 {
		if config, err = config_parser.ParseFileToConfig(config, configpath); err != nil {
			return err
		}
	}
	config_builder := configBuilderWithConfig(config)
	if len(source) > 0 && len(target) > 0 {
		if err := config_builder.AddFile(source, target); err != nil {
			return err
		}
	}
	if len(name) > 0 {
		if err := config_builder.Name(name); err != nil {
			return err
		}
	}
	if len(version) > 0 {
		if err := config_builder.Version(version); err != nil {
			return err
		}
	}
	return package_creator.CreatePackage(config_builder.Build())
}
