# Ditto CLI - Development Standards

## Project Overview
Ditto CLI is a Go-based terminal user interface (TUI) for the Ditto social platform, built using the Bubble Tea framework. It interacts with the Ditto API (defaulting to `http://localhost:9001`).

## Technical Stack
- **Language**: Go 1.23+
- **TUI Framework**: [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Styling**: [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **CLI Framework**: [Cobra](https://github.com/spf13/cobra)
- **API Communication**: Standard `net/http` with custom client in `internal/api`

## Architecture
- `cmd/`: CLI command definitions using Cobra.
- `internal/api/`: API client logic.
- `internal/config/`: Configuration management (auth tokens, base URL, keybindings, theme settings).
- `internal/types/`: Shared data models.
- `internal/ui/theme/`: Central theme system. Defines presets (Tokyo Night, Dracula, etc.) and semantic styles.
- `internal/ui/`: TUI implementation.
    - `model.go`: Root model and state orchestration.
    - `view.go`: Main view rendering and component composition.
    - `keys.go`: Central keybinding definitions.
    - `styles.go`: Type aliases for theme integration.
    - `views/`: Sub-views (feed, login, post detail, command palette, etc.).

## Engineering Mandates
- **Surgical Updates**: When modifying TUI components, ensure minimal disruption to the Bubble Tea update loop.
- **Thematic Consistency**: **Mandatory**: Do NOT use hardcoded `lipgloss.Color` or `lipgloss.NewStyle` with literal colors in views. Use the semantic styles from the active `theme.Theme` (e.g., `m.Theme.TextSubtle`, `m.Theme.AccentText`).
- **Error Handling**: Use the central error reporting system or consistent error wrapping.
- **API Consistency**: All API calls must go through the `internal/api.Client`.
- **Testing**:
    - **Mandatory**: Every new feature or bug fix MUST include corresponding unit tests.
    - **Golden Files**: UI changes affecting rendering must update golden snapshots using `go test ./internal/ui -update`.
    - Add unit tests for new API methods in `internal/api/client_test.go`.
    - Add unit tests for TUI logic in `internal/ui/ui_test.go`.
- **Documentation**:
    - **Mandatory**: Update `README.md` for any changes that affect the user experience (new commands, flags, or configuration).
    - Keep `TASKS.md` updated with progress.
- **Conventions**:
    - Use vim-like keybindings (`j/k` for navigation) where appropriate.
    - Trigger the Command Palette using `:` (colon).
    - Follow standard Go formatting (`go fmt`) and linting rules (`golangci-lint`).
    - Avoid deprecated Lip Gloss methods (e.g., use `style.Padding()` instead of `style.Copy().Padding()`).

## Development Workflow
1. **Build**: `make build`
2. **Test**: `make test`
3. **Lint**: `make lint`
4. **Update Goldens**: `go test ./internal/ui -update`
5. **Run**: `make run` (runs TUI) or `./dist/ditto-cli [command]` for CLI operations.

## Key Files
- `internal/api/client.go`: Main API interface.
- `internal/ui/model.go`: Central state management.
- `internal/ui/theme/theme.go`: Theme presets and resolution logic.
- `cmd/root.go`: CLI entry point.
