package channel

import (
	"errors"
	"testing"
)

func TestTerminalError(t *testing.T) {
	inner := errors.New("account banned")
	term := NewTerminalError(inner)

	if !IsTerminal(term) {
		t.Error("expected IsTerminal to return true")
	}
	if term.Error() != "terminal: account banned" {
		t.Errorf("unexpected error message: %s", term.Error())
	}
	// errors.As should work
	var te *TerminalError
	if !errors.As(term, &te) {
		t.Error("errors.As should unwrap to TerminalError")
	}
	if !te.Terminal() {
		t.Error("Terminal() should return true")
	}
}

func TestIsTerminalNonTerminal(t *testing.T) {
	err := errors.New("regular error")
	if IsTerminal(err) {
		t.Error("regular error should not be terminal")
	}
}

func TestIsTerminalNil(t *testing.T) {
	if IsTerminal(nil) {
		t.Error("nil error should not be terminal")
	}
}
