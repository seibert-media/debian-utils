package zip_extractor

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsZipExtractor(t *testing.T) {
	b := New()
	var i *ZipExtractor
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
