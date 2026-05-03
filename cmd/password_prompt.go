package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"golang.org/x/term"
)

var (
	isTerminal   = term.IsTerminal
	readPassword = term.ReadPassword
)

func resolvePassword(out io.Writer) (string, error) {
	if password != "" {
		return password, nil
	}

	fd := int(os.Stdin.Fd())
	if !isTerminal(fd) {
		return "", fmt.Errorf("password is required")
	}

	if _, err := fmt.Fprint(out, "Password: "); err != nil {
		return "", err
	}
	secret, err := readPassword(fd)
	if _, newlineErr := fmt.Fprintln(out); newlineErr != nil && err == nil {
		err = newlineErr
	}
	if err != nil {
		return "", err
	}

	value := strings.TrimSpace(string(secret))
	if value == "" {
		return "", fmt.Errorf("password is required")
	}
	return value, nil
}
