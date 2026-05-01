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
- `internal/config/`: Configuration management (auth tokens, base URL).
- `internal/types/`: Shared data models.
- `internal/ui/`: TUI implementation.
    - `model.go`: Root model and state.
    - `view.go`: Main view rendering.
    - `styles.go`: Global styles and colors.
    - `views/`: Sub-views (feed, login, post detail, etc.).

## Engineering Mandates
- **Surgical Updates**: When modifying TUI components, ensure minimal disruption to the Bubble Tea update loop.
- **Error Handling**: Use the central error reporting system (to be implemented) or consistent error wrapping.
- **API Consistency**: All API calls must go through the `internal/api.Client`.
- **Testing**:
    - **Mandatory**: Every new feature or bug fix MUST include corresponding unit tests.
    - Add unit tests for new API methods in `internal/api/client_test.go`.
    - Add unit tests for TUI logic in `internal/ui/ui_test.go`, `internal/ui/logic_test.go`, or per-view tests.
- **Documentation**:
    - **Mandatory**: Update `README.md` for any changes that affect the user experience (new commands, flags, or configuration).
    - Keep `TASKS.md` updated with the current progress of the implementation phases.
- **Conventions**:
    - Use vim-like keybindings where appropriate.
    - Follow standard Go formatting (`go fmt`) and linting rules.

## Development Workflow
1. **Build**: `make build`
2. **Test**: `make test`
3. **Lint**: `make lint`
4. **Run**: `make run` (runs TUI) or `./ditto-cli [command]` for CLI operations.

## Key Files
- `internal/api/client.go`: Main API interface.
- `internal/ui/model.go`: Central state management.
- `cmd/root.go`: CLI entry point.
