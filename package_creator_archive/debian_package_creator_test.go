package package_creator_archive

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsDebianPackageCreator(t *testing.T) {
	b := New(nil)
	var i *DebianPackageCreator
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
