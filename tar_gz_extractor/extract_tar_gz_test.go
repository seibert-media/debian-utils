package tar_gz_extractor

import (
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsTarGzExtractor(t *testing.T) {
	b := New()
	var i *TarGzExtractor
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}
