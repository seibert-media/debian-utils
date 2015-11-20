package copier

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsCopier(t *testing.T) {
	b := New()
	var i *Copier
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
