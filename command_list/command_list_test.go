package command_list

import (
	"testing"

	"fmt"

	. "github.com/bborbe/assert"
	"github.com/bborbe/debian-utils/command"
)

func TestImplementsCommandList(t *testing.T) {
	b := New()
	var i *CommandList
	err := AssertThat(b, Implements(i).Message("check type"))
	if err != nil {
		t.Fatal(err)
	}
}

func TestAdd(t *testing.T) {
	list := New()
	if err := AssertThat(len(list.commands), Is(0)); err != nil {
		t.Fatal(err)
	}
	list.Add(command.New(func() error { return nil }, func() error { return nil }))
	if err := AssertThat(len(list.commands), Is(1)); err != nil {
		t.Fatal(err)
	}
	list.Add(command.New(func() error { return nil }, func() error { return nil }))
	if err := AssertThat(len(list.commands), Is(2)); err != nil {
		t.Fatal(err)
	}
}

func TestRunEmpty(t *testing.T) {
	list := New()
	if err := AssertThat(len(list.commands), Is(0)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(list.Run(), NilValue()); err != nil {
		t.Fatal(err)
	}
}

func TestRunSuccess(t *testing.T) {
	list := New()
	doCounter := 0
	undoCounter := 0
	list.Add(command.New(func() error { doCounter++; return nil }, func() error { undoCounter++; return nil }))
	list.Add(command.New(func() error { doCounter++; return nil }, func() error { undoCounter++; return nil }))
	if err := AssertThat(list.Run(), NilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(doCounter, Is(2)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(undoCounter, Is(0)); err != nil {
		t.Fatal(err)
	}
}

func TestRunFirstFail(t *testing.T) {
	list := New()
	doCounter := 0
	undoCounter := 0
	list.Add(command.New(func() error { doCounter++; return fmt.Errorf("foo") }, func() error { undoCounter++; return nil }))
	list.Add(command.New(func() error { doCounter++; return nil }, func() error { undoCounter++; return nil }))
	if err := AssertThat(list.Run(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(doCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(undoCounter, Is(1)); err != nil {
		t.Fatal(err)
	}
}

func TestRunSecondFail(t *testing.T) {
	list := New()
	doCounter := 0
	undoCounter := 0
	list.Add(command.New(func() error { doCounter++; return nil }, func() error { undoCounter++; return nil }))
	list.Add(command.New(func() error { doCounter++; return fmt.Errorf("foo") }, func() error { undoCounter++; return nil }))
	if err := AssertThat(list.Run(), NotNilValue()); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(doCounter, Is(2)); err != nil {
		t.Fatal(err)
	}
	if err := AssertThat(undoCounter, Is(2)); err != nil {
		t.Fatal(err)
	}
}
