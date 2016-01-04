package latest_package_creator

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsLatestDebianPackageCreator(t *testing.T) {
	b := New(nil, nil, nil, nil)
	var i *LatestDebianPackageCreator
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
