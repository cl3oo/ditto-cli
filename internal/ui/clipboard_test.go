package ui

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestClipboardCommands(t *testing.T) {
	t.Run("darwin", func(t *testing.T) {
		cmds := clipboardCommands("darwin")
		if len(cmds) != 1 || cmds[0].name != "pbcopy" {
			t.Fatalf("unexpected darwin commands: %#v", cmds)
		}
	})

	t.Run("linux without wayland", func(t *testing.T) {
		old := os.Getenv("WAYLAND_DISPLAY")
		t.Cleanup(func() { _ = os.Setenv("WAYLAND_DISPLAY", old) })
		_ = os.Unsetenv("WAYLAND_DISPLAY")

		cmds := clipboardCommands("linux")
		if len(cmds) < 3 {
			t.Fatalf("expected linux fallbacks, got %#v", cmds)
		}
		if cmds[0].name != "xclip" {
			t.Fatalf("expected xclip first without wayland, got %#v", cmds)
		}
	})

	t.Run("linux with wayland", func(t *testing.T) {
		old := os.Getenv("WAYLAND_DISPLAY")
		t.Cleanup(func() { _ = os.Setenv("WAYLAND_DISPLAY", old) })
		if err := os.Setenv("WAYLAND_DISPLAY", "wayland-0"); err != nil {
			t.Fatalf("set env: %v", err)
		}

		cmds := clipboardCommands("linux")
		if len(cmds) == 0 || cmds[0].name != "wl-copy" {
			t.Fatalf("expected wl-copy first with wayland, got %#v", cmds)
		}
	})
}

func TestPipeToCommandWritesViaStdin(t *testing.T) {
	tmp := t.TempDir()
	output := filepath.Join(tmp, "clipboard.txt")
	script := filepath.Join(tmp, "fake-copy.sh")
	scriptBody := "#!/bin/sh\ncat > \"$TEST_OUTPUT\"\n"
	if err := os.WriteFile(script, []byte(scriptBody), 0o755); err != nil {
		t.Fatalf("write script: %v", err)
	}

	oldExec := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		cmd := oldExec(script)
		cmd.Env = append(os.Environ(), "TEST_OUTPUT="+output)
		return cmd
	}
	t.Cleanup(func() { execCommand = oldExec })

	if err := pipeToCommand("https://ditto.example/post/123\n", "ignored"); err != nil {
		t.Fatalf("pipeToCommand failed: %v", err)
	}

	data, err := os.ReadFile(output)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if got := string(data); got != "https://ditto.example/post/123\n" {
		t.Fatalf("unexpected clipboard payload: %q", got)
	}
}

func TestCopyToClipboardFallsBackWhenCommandMissing(t *testing.T) {
	oldExec := execCommand
	execCommand = func(name string, args ...string) *exec.Cmd {
		return oldExec("/usr/bin/env", append([]string{"__definitely_missing_command__"}, args...)...)
	}
	t.Cleanup(func() { execCommand = oldExec })

	if err := copyToClipboard("hello"); err == nil || !strings.Contains(err.Error(), "no clipboard command available") {
		t.Fatalf("expected fallback error, got %v", err)
	}
}
