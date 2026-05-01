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

- [Go](https://go.dev/doc/install) 1.26.2 or higher (matches `go.mod`).

### From Source

```bash
# Clone the repository
git clone https://github.com/rfcku/ditto-cli.git
cd ditto-cli

# Build the binary
make build

# Launch the TUI
./ditto-cli tui
```

### Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `DITTO_API_URL` | The URL of the Ditto API. | `http://localhost:9001/v1` |

## Usage

### TUI Commands

While in the TUI, press `:` to enter command mode:

| Command | Action |
| :--- | :--- |
| `:f`, `:feed` | Return to the home feed |
| `:c`, `:communities` | View communities |
| `:n`, `:new` | Create a new post |
| `:s <query>` | Search communities |
| `:edit` | Edit the current post or community |
| `:delete` | Delete the current post or community |
| `:random` | Fetch 10 random posts or communities |
| `:follow` | Follow/Unfollow the current post author |
| `:following` | List users you follow |
| `:joined` | List communities you joined |
| `:ban <userID>` | Ban a user from the community (Moderator only) |
| `:mod add <userID>` | Add a moderator to the community (Moderator only) |
| `:lock` | Lock the current post (Moderator only) |
| `:unlock` | Unlock the current post (Moderator only) |
| `:mod-delete <reason>` | Delete a post as a moderator |
| `:award <award_id>` | Give an award to the current post |
| `:report <reason>` | Report the current post or community |
| `:delete-comment <id>` | Delete a specific comment by ID |
| `:settings` | Edit your user profile (avatar) |
| `:delete-account` | Permanently delete your account |
| `:q`, `:quit` | Exit the application |

### CLI Mode

You can also use `ditto-cli` for one-off operations:

```bash
# Launch the TUI without building first
go run . tui

# Upload media to a post
./ditto-cli upload image.png --target <post_id>

# Download media
./ditto-cli download <media_id> --output photo.png

# Get random posts, communities, users or comments
./ditto-cli get posts --random 5
./ditto-cli get communities --random 10
./ditto-cli get users --random 3
./ditto-cli get comments --random 5
```

## Development

### Key Commands

- `make run`: Run the TUI in development mode (`go run . tui`).
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

There is currently no GitHub Actions workflow checked into this repository, so local `make test` and `make lint` are the primary verification steps for now.
