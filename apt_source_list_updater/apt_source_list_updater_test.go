package apt_source_list_updater

import (
	"testing"

	"runtime"

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

func TestParseLine(t *testing.T) {
	var err error
	var infos *infos
	infos, err = ParseLine("deb https://example.com/repo dist comp")
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.url, Is("https://example.com/repo")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.distribution, Is("dist")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.component, Is("comp")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.architecture, Is(runtime.GOARCH)); err != nil {
		t.Fatal(err)
	}
}

func TestParseLineWithArch(t *testing.T) {
	var err error
	var infos *infos
	infos, err = ParseLine("deb [arch=all] https://example.com/repo dist comp")
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.url, Is("https://example.com/repo")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.distribution, Is("dist")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.component, Is("comp")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.architecture, Is("all")); err != nil {
		t.Fatal(err)
	}
}

func TestParseLineWithAmd64(t *testing.T) {
	var err error
	var infos *infos
	infos, err = ParseLine("deb [arch=amd64] http://aptly.benjamin-borbe.de/stable  default main\n")
	if err = AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.url, Is("http://aptly.benjamin-borbe.de/stable")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.distribution, Is("default")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.component, Is("main")); err != nil {
		t.Fatal(err)
	}
	if err = AssertThat(infos.architecture, Is("amd64")); err != nil {
		t.Fatal(err)
	}
}

func TestRemotePackagesUrl(t *testing.T) {
	infos := &infos{
		url:          "http://www.example.com/repo",
		distribution: "default",
		architecture: "all",
		component:    "main",
	}
	if err := AssertThat(infos.RemotePackagesUrl(), Is("http://www.example.com/repo/dists/default/main/binary-all/Packages")); err != nil {
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
