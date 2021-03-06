package main

import (
	"testing"

	. "github.com/bborbe/assert"
	debian_command_list "github.com/bborbe/command/list"
	debian_config "github.com/seibert-media/debian-utils/config"
	debian_config_builder "github.com/seibert-media/debian-utils/config_builder"
	debian_config_parser "github.com/seibert-media/debian-utils/config_parser"
	debian_copier "github.com/seibert-media/debian-utils/copier"
	debian_package_creator "github.com/seibert-media/debian-utils/package_creator"
)

func TestDo(t *testing.T) {
	var err error

	commandProvider := func() debian_command_list.CommandList {
		return debian_command_list.New()
	}
	configBuilderWithConfig := func(config *debian_config.Config) debian_config_builder.ConfigBuilder {
		return debian_config_builder.NewWithConfig(config)
	}
	copier := debian_copier.New()
	package_creator := debian_package_creator.New(commandProvider, copier, nil, nil, nil, nil)
	config_parser := debian_config_parser.New()

	err = do(config_parser, configBuilderWithConfig, package_creator, "", "", "", "", "")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
