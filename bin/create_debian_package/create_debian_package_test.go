package main

import (
	"testing"

	. "github.com/bborbe/assert"
	debian_command_list "github.com/bborbe/debian-utils/command_list"
	debian_config_builder "github.com/bborbe/debian-utils/config_builder"
	debian_package_creator "github.com/bborbe/debian-utils/package_creator"
	io_mock "github.com/bborbe/io/mock"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()
	config_builder := debian_config_builder.New()
	package_creator := debian_package_creator.New(func() debian_command_list.CommandList {
		return debian_command_list.New()
	}, nil)

	err = do(writer, config_builder, package_creator, "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
