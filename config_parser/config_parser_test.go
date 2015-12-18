package config_parser

import (
	"testing"

	debian_config "github.com/bborbe/debian_utils/config"

	. "github.com/bborbe/assert"
)

func TestImplementsConfigParser(t *testing.T) {
	b := New()
	var i *ConfigParser
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestDefaults(t *testing.T) {
	config := debian_config.DefaultConfig()
	config, err := New().ParseContentToConfig(config, []byte(`{}`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config.Section, Is("base")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config.Priority, Is("optional")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config.Architecture, Is("all")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config.Maintainer, Is("Benjamin Borbe <bborbe@rocketnews.de>")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config.Description, Is("-")); err != nil {
		t.Fatal(err)
	}
}

func TestParseConfigArchitecture(t *testing.T) {
	config := debian_config.DefaultConfig()
	config, err := New().ParseContentToConfig(config, []byte(`{"architecture":"amd64"}`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config.Architecture, Is("amd64")); err != nil {
		t.Fatal(err)
	}
}

func TestParseConfigName(t *testing.T) {
	config := debian_config.DefaultConfig()
	config, err := New().ParseContentToConfig(config, []byte(`{"name":"helloworld"}`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config.Name, Is("helloworld")); err != nil {
		t.Fatal(err)
	}
}

func TestParseConfigVersion(t *testing.T) {
	config := debian_config.DefaultConfig()
	config, err := New().ParseContentToConfig(config, []byte(`{"version":"1.2.3"}`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config.Version, Is("1.2.3")); err != nil {
		t.Fatal(err)
	}
}

func TestParseConfigFilesEmpty(t *testing.T) {
	config := debian_config.DefaultConfig()
	config, err := New().ParseContentToConfig(config, []byte(`{"files":[]}`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(config.Files), Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestParseConfigFiles(t *testing.T) {
	config := debian_config.DefaultConfig()
	config, err := New().ParseContentToConfig(config, []byte(`{"files":[{"source":"/tmp/source.txt","target":"/tmp/target.txt"}]}`))
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(len(config.Files), Is(1)); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config.Files[0].Source, Is("/tmp/source.txt")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(config.Files[0].Target, Is("/tmp/target.txt")); err != nil {
		t.Fatal(err)
	}
}
