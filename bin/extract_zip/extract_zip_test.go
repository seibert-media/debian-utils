package main

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestDo(t *testing.T) {
	err := do(nil, "", "")
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
}
