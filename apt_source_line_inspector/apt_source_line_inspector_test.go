package line_inspector

import (
	"testing"

	"runtime"

	. "github.com/bborbe/assert"
)

func TestImplementsAptSourceListUpdater(t *testing.T) {
	b := New(nil)
	var i *LineInspector
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestParseLine(t *testing.T) {
	var err error
	var infos *infos
	infos, err = parseLine("deb https://example.com/repo dist comp")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.url, Is("https://example.com/repo")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.distribution, Is("dist")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.component, Is("comp")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.architecture, Is(runtime.GOARCH)); err != nil {
		t.Fatal(err)
	}
}

func TestParseLineWithArch(t *testing.T) {
	var err error
	var infos *infos
	infos, err = parseLine("deb [arch=all] https://example.com/repo dist comp")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.url, Is("https://example.com/repo")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.distribution, Is("dist")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.component, Is("comp")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.architecture, Is("all")); err != nil {
		t.Fatal(err)
	}
}

func TestParseLineWithAmd64(t *testing.T) {
	var err error
	var infos *infos
	infos, err = parseLine("deb [arch=amd64] https://www.benjamin-borbe.de/aptly/stable  default main\n")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.url, Is("https://www.benjamin-borbe.de/aptly/stable")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.distribution, Is("default")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.component, Is("main")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.architecture, Is(runtime.GOARCH)); err != nil {
		t.Fatal(err)
	}
}

func TestParseLineWithoutComp(t *testing.T) {
	var err error
	var infos *infos
	infos, err = parseLine("deb https://www.benjamin-borbe.de/aptly/stable default\n")
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.url, Is("https://www.benjamin-borbe.de/aptly/stable")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.distribution, Is("default")); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(infos.architecture, Is("amd64")); err != nil {
		t.Fatal(err)
	}
}

func TestRemotePackagesURL(t *testing.T) {
	infos := &infos{
		url:          "http://www.example.com/repo",
		distribution: "default",
		architecture: "all",
		component:    "main",
	}
	if err := AssertThat(infos.RemotePackagesURL(), Is("http://www.example.com/repo/dists/default/main/binary-all/Packages")); err != nil {
		t.Fatal(err)
	}
}

func TestLocalPackagesFile(t *testing.T) {
	infos := &infos{
		url:          "http://www.example.com/repo",
		distribution: "default",
		architecture: "all",
		component:    "main",
	}
	if err := AssertThat(infos.LocalPackagesFile(), Is("/var/lib/apt/lists/www.example.com_repo_dists_default_main_binary-all_Packages")); err != nil {
		t.Fatal(err)
	}
}

func TestLocalPackagesFileWithAuth(t *testing.T) {
	infos := &infos{
		url:          "http://user:passs@www.example.com/repo",
		distribution: "default",
		architecture: "all",
		component:    "main",
	}
	if err := AssertThat(infos.LocalPackagesFile(), Is("/var/lib/apt/lists/www.example.com_repo_dists_default_main_binary-all_Packages")); err != nil {
		t.Fatal(err)
	}
}
