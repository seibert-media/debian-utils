package main

import (
	"testing"

	"bytes"

	. "github.com/bborbe/assert"
	debian_command_list "github.com/bborbe/command/list"
	debian_config "github.com/bborbe/debian_utils/config"
	debian_config_builder "github.com/bborbe/debian_utils/config_builder"
	debian_config_parser "github.com/bborbe/debian_utils/config_parser"
	debian_copier "github.com/bborbe/debian_utils/copier"
	debian_package_creator "github.com/bborbe/debian_utils/package_creator"
)

func TestDo(t *testing.T) {
	var err error
	writer := bytes.NewBufferString("")

	commandProvider := func() debian_command_list.CommandList {
		return debian_command_list.New()
	}
	configBuilderWithConfig := func(config *debian_config.Config) debian_config_builder.ConfigBuilder {
		return debian_config_builder.NewWithConfig(config)
	}
	copier := debian_copier.New()
	package_creator := debian_package_creator.New(commandProvider, copier)
	config_parser := debian_config_parser.New()

	err = do(writer, config_parser, configBuilderWithConfig, package_creator, "", "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
