package package_creator


import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsImageSaver(t *testing.T) {
	b := New()
	var i *Builder
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}