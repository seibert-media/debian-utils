package main

import (
	"testing"

	. "github.com/bborbe/assert"
	io_mock "github.com/bborbe/io/mock"
	"github.com/bborbe/debian/package_builder"
	"github.com/bborbe/debian/command_list"
)

func TestDo(t *testing.T) {
	var err error
	writer := io_mock.NewWriter()
	builder := package_creator.New(command_list.New())
	err = do(writer, builder, "", "", "", "")
	err = AssertThat(err, NotNilValue())
	if err != nil {
		t.Fatal(err)
	}
}
