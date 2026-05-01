# Ditto CLI ⌨️

A powerful terminal-based user interface (TUI) for the Ditto social platform. Built with Go and the [Bubble Tea](https://github.com/charmbracelet/bubbletea) framework.

## Features

- 🏠 **Home Feed**: View trending and personalized posts.
- 🗳 **Voting**: Upvote and downvote content directly from the terminal.
- 💬 **Interactions**: Read and create comments.
- 🏘 **Communities**: Browse and join communities.
- 👤 **Profile**: View your wallet balance and user details.
- ⌨️ **Keyboard Optimized**: Fast navigation with vim-like bindings or arrows.
- 🧭 **Contextual Help Bar**: A sticky footer advertises the most useful actions for the current screen.

## Installation

### Prerequisites

- [Go](https://go.dev/doc/install) 1.26.2 or higher (matches `go.mod`).
- A reachable Ditto API, local or remote.

### Quick start from source

```bash
git clone https://github.com/rfcku/ditto-cli.git
cd ditto-cli
make build
./dist/ditto-cli tui
```

### Install the CLI with Go

```bash
go install github.com/rfcku/ditto-cli@latest

# Then run either mode
ditto-cli tui
ditto-cli get posts --random 5
```

### Install prebuilt binaries from GitHub Releases

Download the archive that matches your OS and CPU from the Releases page, then extract it and move `ditto-cli` somewhere on your `PATH`.

```bash
# Example for macOS Apple Silicon
curl -L https://github.com/rfcku/ditto-cli/releases/latest/download/ditto-cli_VERSION_Darwin_arm64.tar.gz | tar -xz
chmod +x ditto-cli
mv ditto-cli /usr/local/bin/
```

### Environment Variables

| Variable | Description | Default |
| :--- | :--- | :--- |
| `DITTO_API_URL` | Base API URL used by the CLI and TUI. | `http://localhost:9001/api/v1` |

## Usage

Ditto CLI supports two workflows. You can also check the exact build metadata with `ditto-cli version`:

- **TUI mode** for browsing, posting, moderation, and day-to-day navigation inside the terminal.
- **CLI mode** for one-off actions, scripting, smoke tests, and quick API checks.

### TUI Commands

While in the TUI, press `:` to enter command mode:

| Command | Action |
| :--- | :--- |
| `:f`, `:feed` | Return to the home feed |
| `:c`, `:communities` | View communities |
| `:n`, `:new` | Create a new post |
| `:s <query>` | Search communities and posts |
| `:man` | Open the user manual |
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

### Navigation & Shortcuts

The TUI also keeps a sticky footer at the bottom with context-aware hints, so the current screen tells you what matters without opening the manual.

- **j / k** or **Up / Down**: Navigate lists.
- **Enter**: Select or submit.
- **q**: Go back.
- **a**: Upvote (configurable).
- **z**: Downvote (configurable).
- **r**: Reply to a comment (when selected).
- **p**: View user profile (when comment selected).
- **L**: Load more comments.
- **R**: Report post or comment.
- **s**: Share post link.

### CLI Mode

You can also use `ditto-cli` for one-off operations:

```bash
# Launch the TUI without building first
go run . tui

# Upload media to a post
./dist/ditto-cli upload image.png --target <post_id>

# Download media
./dist/ditto-cli download <media_id> --output photo.png

# Get random posts, communities, users or comments
./dist/ditto-cli get posts --random 5
./dist/ditto-cli get communities --random 10
./dist/ditto-cli get users --random 3
./dist/ditto-cli get comments --random 5
```

## Development

### Contributor workflow

```bash
git clone https://github.com/rfcku/ditto-cli.git
cd ditto-cli
make test
make lint
make build
```

### Key Commands

- `make run`: Run the TUI in development mode (`go run . tui`).
- `make test`: Run unit tests.
- `make lint`: Run the linter (`golangci-lint`).
- `make build`: Compile the binary into `dist/ditto-cli`.
- `make smoke`: Run a product smoke test for register, community creation, post creation, feed reads, join, and vote flows against a live Ditto API.
- `make release-check`: Validate the GoReleaser config.
- `make release-build`: Build release archives locally without publishing.
- `make release-test`: Run a full local snapshot release without publishing.
- `make clean`: Remove generated build artifacts from `dist/`.

### Smoke Testing Against a Live API

`make smoke` uses `scripts/smoke.sh` and expects a reachable Ditto API. By default it targets `http://localhost:9001/api/v1`, creates two disposable users in an isolated temporary config directory, then exercises the core CLI journey end-to-end.

```bash
# Use the default local API
make smoke

# Or point the smoke test at another environment
DITTO_API_URL=https://ditto.example.com/api/v1 make smoke
```

Useful overrides:

- `DITTO_SMOKE_BIN` to point at a custom built binary.
- `DITTO_SMOKE_USER_ONE`, `DITTO_SMOKE_USER_TWO`, `DITTO_SMOKE_EMAIL_ONE`, `DITTO_SMOKE_EMAIL_TWO` if you need deterministic test identities.

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

GitHub Actions runs test, lint, build, and release-config validation checks on pushes to `main` and on pull requests. Local `make test`, `make lint`, `make build`, and `make release-check` should match that baseline before opening a PR. For live-environment validation, `make smoke` provides a repeatable product-level check outside the default CI matrix.

Tagged releases (`v*`) publish prebuilt archives plus checksums for macOS and Linux on `amd64` and `arm64` through GoReleaser.
