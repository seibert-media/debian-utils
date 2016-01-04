package package_creator_by_reader

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsDebianPackageCreator(t *testing.T) {
	b := New(nil, nil)
	var i *Creator
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestJoinDirsOne(t *testing.T) {
	if err := AssertThat(joinDirs("/tmp"), Is("/tmp")); err != nil {
		t.Fatal(err)
	}
}

func TestJoinDirsTwo(t *testing.T) {
	if err := AssertThat(joinDirs("/tmp", "sub"), Is("/tmp/sub")); err != nil {
		t.Fatal(err)
	}
}

func TestJoinDirsOneEmpty(t *testing.T) {
	if err := AssertThat(joinDirs("/tmp", ""), Is("/tmp")); err != nil {
		t.Fatal(err)
	}
}