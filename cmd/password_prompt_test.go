package cmd

import (
	"bytes"
	"errors"
	"testing"
)

func TestResolvePassword_UsesFlagValue(t *testing.T) {
	prevPassword := password
	defer func() { password = prevPassword }()

	password = "secret"
	got, err := resolvePassword(&bytes.Buffer{})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "secret" {
		t.Fatalf("expected secret, got %q", got)
	}
}

func TestResolvePassword_PromptsWhenInteractive(t *testing.T) {
	prevPassword := password
	prevTerminal := isTerminal
	prevReader := readPassword
	defer func() {
		password = prevPassword
		isTerminal = prevTerminal
		readPassword = prevReader
	}()

	password = ""
	isTerminal = func(int) bool { return true }
	readPassword = func(int) ([]byte, error) { return []byte("typed-secret\n"), nil }

	buf := &bytes.Buffer{}
	got, err := resolvePassword(buf)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got != "typed-secret" {
		t.Fatalf("expected typed-secret, got %q", got)
	}
	if buf.String() != "Password: \n" {
		t.Fatalf("unexpected prompt output %q", buf.String())
	}
}

func TestResolvePassword_RequiresPasswordWhenNonInteractive(t *testing.T) {
	prevPassword := password
	prevTerminal := isTerminal
	defer func() {
		password = prevPassword
		isTerminal = prevTerminal
	}()

	password = ""
	isTerminal = func(int) bool { return false }

	_, err := resolvePassword(&bytes.Buffer{})
	if err == nil || err.Error() != "password is required" {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestResolvePassword_PropagatesReadError(t *testing.T) {
	prevPassword := password
	prevTerminal := isTerminal
	prevReader := readPassword
	defer func() {
		password = prevPassword
		isTerminal = prevTerminal
		readPassword = prevReader
	}()

	password = ""
	isTerminal = func(int) bool { return true }
	readPassword = func(int) ([]byte, error) { return nil, errors.New("boom") }

	_, err := resolvePassword(&bytes.Buffer{})
	if err == nil || err.Error() != "boom" {
		t.Fatalf("unexpected error: %v", err)
	}
}
