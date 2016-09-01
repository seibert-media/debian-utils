package package_creator

import (
	"testing"

	. "github.com/bborbe/assert"
	debian_config "github.com/bborbe/debian_utils/config"
)

func TestImplementsPackageCreator(t *testing.T) {
	b := New(nil, nil, nil, nil, nil, nil)
	var i *PackageCreator
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestControlContentDefault(t *testing.T) {
	config := debian_config.DefaultConfig()
	config.Name = "testPackage"
	config.Section = "utils"
	config.Description = "my test package"
	config.Version = "1.2.3"
	content := controlContent(*config)
	if err := AssertThat(len(content), Not(Eq(0))); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(string(content), Is(`Package: testPackage
Version: 1.2.3
Section: utils
Priority: optional
Architecture: all
Maintainer: Benjamin Borbe <bborbe@rocketnews.de>
Description: my test package
`)); err != nil {
		t.Fatal(err)
	}
}

func TestControlContentDepends(t *testing.T) {
	{
		config := debian_config.DefaultConfig()
		config.Depends = []string{"a"}
		content := controlContent(*config)
		if err := AssertThat(len(content), Not(Eq(0))); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(string(content), Contains("\nDepends: a\n")); err != nil {
			t.Fatal(err)
		}
	}
	{
		config := debian_config.DefaultConfig()
		config.Depends = []string{"a", "b", "c"}
		content := controlContent(*config)
		if err := AssertThat(len(content), Not(Eq(0))); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(string(content), Contains("\nDepends: a,b,c\n")); err != nil {
			t.Fatal(err)
		}
	}
}

func TestControlContentReplaces(t *testing.T) {
	{
		config := debian_config.DefaultConfig()
		config.Replaces = []string{"a"}
		content := controlContent(*config)
		if err := AssertThat(len(content), Not(Eq(0))); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(string(content), Contains("\nReplaces: a\n")); err != nil {
			t.Fatal(err)
		}
	}
	{
		config := debian_config.DefaultConfig()
		config.Replaces = []string{"a", "b", "c"}
		content := controlContent(*config)
		if err := AssertThat(len(content), Not(Eq(0))); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(string(content), Contains("\nReplaces: a,b,c\n")); err != nil {
			t.Fatal(err)
		}
	}
}

func TestControlContentProvides(t *testing.T) {
	{
		config := debian_config.DefaultConfig()
		config.Provides = []string{"a"}
		content := controlContent(*config)
		if err := AssertThat(len(content), Not(Eq(0))); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(string(content), Contains("\nProvides: a\n")); err != nil {
			t.Fatal(err)
		}
	}
	{
		config := debian_config.DefaultConfig()
		config.Provides = []string{"a", "b", "c"}
		content := controlContent(*config)
		if err := AssertThat(len(content), Not(Eq(0))); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(string(content), Contains("\nProvides: a,b,c\n")); err != nil {
			t.Fatal(err)
		}
	}
}

func TestControlContentConflicts(t *testing.T) {
	{
		config := debian_config.DefaultConfig()
		config.Conflicts = []string{"a"}
		content := controlContent(*config)
		if err := AssertThat(len(content), Not(Eq(0))); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(string(content), Contains("\nConflicts: a\n")); err != nil {
			t.Fatal(err)
		}
	}
	{
		config := debian_config.DefaultConfig()
		config.Conflicts = []string{"a", "b", "c"}
		content := controlContent(*config)
		if err := AssertThat(len(content), Not(Eq(0))); err != nil {
			t.Fatal(err)
		}
		if err := AssertThat(string(content), Contains("\nConflicts: a,b,c\n")); err != nil {
			t.Fatal(err)
		}
	}
}
