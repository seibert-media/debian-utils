package url_downloader

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsUrlDownloader(t *testing.T) {
	b := New(nil, nil)
	var i *UrlDownloader
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
