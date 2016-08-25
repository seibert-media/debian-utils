package url_downloader

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsURLDownloader(t *testing.T) {
	b := New(nil, nil)
	var i *URLDownloader
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
