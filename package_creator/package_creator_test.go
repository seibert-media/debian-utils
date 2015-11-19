package package_creator

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsPackageCreator(t *testing.T) {
	b := New(nil)
	var i *PackageCreator
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
