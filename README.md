# Ditto CLI ⌨️

A powerful terminal-based user interface (TUI) for the Ditto social platform. Built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

## Features

- 🏠 **Home Feed**: View trending and personalized posts.
- 🗳 **Voting**: Upvote and downvote content directly from the terminal.
- 💬 **Interactions**: Read and create comments.
- 🏘 **Communities**: Browse and join communities.
- 👤 **Profile**: View your wallet balance and user details.
- ⌨️ **Keyboard Optimized**: Fast navigation with vim-like bindings or arrows.

## Installation

### Prerequisites

- [Go](https://go.dev/doc/install) 1.23 or higher.

### From Source

```bash
# Clone the repository (if not already done)
git clone https://github.com/rfcku/ditto.git
cd ditto/cli

# Build the binary
make build

# Run directly
./ditto-cli
```

### Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `DITTO_API_URL` | The URL of the Ditto API. | `http://localhost:9001` |

## Development

### Key Commands

- `make run`: Run the TUI in development mode.
- `make test`: Run unit tests.
- `make lint`: Run the linter (`golangci-lint`).
- `make build`: Compile the binary.

### Quality Control

This project uses git hooks to ensure code quality:
- **Pre-commit**: Runs `make lint`.
- **Pre-push**: Runs `make test`.

### Architecture

- `internal/api/`: Handles all HTTP communication with the Ditto backend.
- `internal/ui/`: Contains the TUI logic using Bubble Tea.
    - `model.go`: State management.
    - `view.go`: Lip Gloss styling and layout.
    - `update.go`: Event handling.

## CI/CD

The CLI is automatically tested and linted via GitHub Actions on every push to the `cli/` directory.
