package apt_source_has_changed

import (
	"fmt"
	"io"
	"testing"

	. "github.com/bborbe/assert"
)

func TestImplementsAptSourceListUpdater(t *testing.T) {
	b := New(nil)
	var i *AptSourceHasChanged
	if err := AssertThat(b, Implements(i).Message("check type")); err != nil {
		t.Fatal(err)
	}
}

// err == eof ==> return haschangedresult
// err == nil ==> return haschangedresult
// err != nil && err != eof => return false,err

func TestHasFileChangedReturnError(t *testing.T) {
	hasLineChangedCounter := 0
	readStringCounter := 0
	hasLineChanged := func(line string) (bool, error) {
		hasLineChangedCounter++
		return false, nil
	}
	readString := func(delim byte) (line string, err error) {
		readStringCounter++
		return "foo", fmt.Errorf("custom error")
	}
	result, err := hasFileChanged(readString, hasLineChanged)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result, Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(readStringCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hasLineChangedCounter, Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestHasFileChangedReturnEofFalse(t *testing.T) {
	hasLineChangedCounter := 0
	readStringCounter := 0
	hasLineChanged := func(line string) (bool, error) {
		hasLineChangedCounter++
		return false, nil
	}
	readString := func(delim byte) (line string, err error) {
		readStringCounter++
		return "foo", io.EOF
	}
	result, err := hasFileChanged(readString, hasLineChanged)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result, Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(readStringCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hasLineChangedCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestHasFileChangedReturnEofTrue(t *testing.T) {
	hasLineChangedCounter := 0
	readStringCounter := 0
	hasLineChanged := func(line string) (bool, error) {
		hasLineChangedCounter++
		return true, nil
	}
	readString := func(delim byte) (line string, err error) {
		readStringCounter++
		return "foo", io.EOF
	}
	result, err := hasFileChanged(readString, hasLineChanged)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result, Is(true)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(readStringCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hasLineChangedCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestHasFileChangedReturnEofError(t *testing.T) {
	hasLineChangedCounter := 0
	readStringCounter := 0
	hasLineChanged := func(line string) (bool, error) {
		hasLineChangedCounter++
		return false, fmt.Errorf("foo")
	}
	readString := func(delim byte) (line string, err error) {
		readStringCounter++
		return "foo", io.EOF
	}
	result, err := hasFileChanged(readString, hasLineChanged)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result, Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(readStringCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hasLineChangedCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestHasFileChangedTrue(t *testing.T) {
	hasLineChangedCounter := 0
	readStringCounter := 0
	hasLineChanged := func(line string) (bool, error) {
		hasLineChangedCounter++
		return true, nil
	}
	readString := func(delim byte) (line string, err error) {
		readStringCounter++
		return "foo", nil
	}
	result, err := hasFileChanged(readString, hasLineChanged)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result, Is(true)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(readStringCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hasLineChangedCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestHasFileChangedFalse(t *testing.T) {
	hasLineChangedCounter := 0
	readStringCounter := 0
	hasLineChanged := func(line string) (bool, error) {
		hasLineChangedCounter++
		return false, nil
	}
	readString := func(delim byte) (line string, err error) {
		readStringCounter++
		if readStringCounter > 1 {
			return "", io.EOF
		}
		return "foo", nil
	}
	result, err := hasFileChanged(readString, hasLineChanged)
	if err := AssertThat(err, NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result, Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(readStringCounter, Is(2)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hasLineChangedCounter, Is(2)); err != nil {
		t.Fatal(err)
	}
}

func TestHasFileChangedError(t *testing.T) {
	hasLineChangedCounter := 0
	readStringCounter := 0
	hasLineChanged := func(line string) (bool, error) {
		hasLineChangedCounter++
		return false, fmt.Errorf("custom error")
	}
	readString := func(delim byte) (line string, err error) {
		readStringCounter++
		return "foo", nil
	}
	result, err := hasFileChanged(readString, hasLineChanged)
	if err := AssertThat(err, NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(result, Is(false)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(readStringCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(hasLineChangedCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
}
