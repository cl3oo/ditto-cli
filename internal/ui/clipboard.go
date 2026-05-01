package ui

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"runtime"
)

var execCommand = exec.Command

func copyToClipboard(text string) error {
	commands := clipboardCommands(runtime.GOOS)
	if len(commands) == 0 {
		return fmt.Errorf("clipboard unsupported on %s", runtime.GOOS)
	}

	for _, cmd := range commands {
		if err := pipeToCommand(text, cmd.name, cmd.args...); err == nil {
			return nil
		}
	}

	return fmt.Errorf("no clipboard command available")
}

type clipboardCommand struct {
	name string
	args []string
}

func clipboardCommands(goos string) []clipboardCommand {
	switch goos {
	case "darwin":
		return []clipboardCommand{{name: "pbcopy"}}
	case "linux":
		commands := []clipboardCommand{}
		if os.Getenv("WAYLAND_DISPLAY") != "" {
			commands = append(commands, clipboardCommand{name: "wl-copy"})
		}
		commands = append(commands,
			clipboardCommand{name: "xclip", args: []string{"-selection", "clipboard"}},
			clipboardCommand{name: "wl-copy"},
			clipboardCommand{name: "xsel", args: []string{"--clipboard", "--input"}},
		)
		return commands
	default:
		return nil
	}
}

func pipeToCommand(text, name string, args ...string) error {
	cmd := execCommand(name, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		_ = stdin.Close()
		return err
	}

	_, writeErr := io.WriteString(stdin, text)
	closeErr := stdin.Close()
	waitErr := cmd.Wait()

	if writeErr != nil {
		return writeErr
	}
	if closeErr != nil {
		return closeErr
	}
	return waitErr
}
