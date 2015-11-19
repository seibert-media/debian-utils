package config_builder

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsConfigBuilder(t *testing.T) {
	b := New()
	var i *ConfigBuilder
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
