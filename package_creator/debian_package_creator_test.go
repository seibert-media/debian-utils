package package_creator


import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsImageSaver(t *testing.T) {
	c := New()
	var i *Creator
	err := AssertThat(c, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}