package cmd

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelpClarifiesTUIEntryPoint(t *testing.T) {
	buf := &bytes.Buffer{}
	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)
	rootCmd.SetArgs([]string{"--help"})
	defer rootCmd.SetArgs(nil)

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("root help failed: %v", err)
	}

	output := buf.String()
	for _, want := range []string{
		"ditto-cli tui",
		"run 'ditto-cli <command>' for one-off CLI actions",
		"ditto-cli get posts --random 5",
	} {
		if !strings.Contains(output, want) {
			t.Fatalf("help output %q missing %q", output, want)
		}
	}
}
