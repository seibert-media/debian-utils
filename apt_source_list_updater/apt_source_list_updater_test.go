package apt_source_list_updater

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsAptSourceListUpdater(t *testing.T) {
	b := New(nil)
	var i *AptSourceListUpdater
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
