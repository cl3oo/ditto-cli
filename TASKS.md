# Ditto CLI Development Roadmap 🚀

This document outlines the professional engineering tasks required to transform the Ditto CLI into a fully integrated, robust, and well-tested TUI application connected to `api.ditto.local`.

## Phase 1: Core API & Auth Stability 🔐
- [ ] **Persistent Session Management**: Implement secure token storage (e.g., using `keyring` or local encrypted config) to avoid repeated logins.
- [ ] **Global Error Handling**: Create a unified error reporting system in the TUI (e.g., a "toast" notification or error modal).
- [ ] **API Interceptors**: Add middleware to the `Client` to automatically attach JWT tokens and handle `401 Unauthorized` by redirecting to login.
- [ ] **Request Logging**: Implement a debug log file (e.g., `ditto.log`) to trace API requests and responses during development.

## Phase 2: Feature Integration (UI & Logic) 🛠
- [ ] **Infinite Scrolling (Feed)**: Implement pagination logic in the `FeedModel` to fetch more posts as the user scrolls to the bottom.
- [ ] **Interactive Comments**: 
    - [ ] Add ability to "Reply" to comments within the `PostDetailView`.
    - [ ] Implement voting (up/down) for individual comments.
- [ ] **Voting System**:
    - [ ] Add visual feedback for votes (e.g., color changes for upvoted items).
    - [ ] Implement "undo" vote capability.
- [ ] **Community Management**:
    - [ ] Implement "Join/Leave" functionality in the Communities view.
    - [ ] Add "Create Community" modal.
- [ ] **Rich Media Support**: Add a "View Image/Link" action that opens the system's default browser/viewer.
- [ ] **Search**: Implement a global search feature for posts, users, and communities.

## Phase 3: UX & Performance 🎨
- [ ] **Consistent Navigation**: Standardize "Esc" and "Backspace" for navigation history (backstack).
- [ ] **Loading States**: Add a spinner component (`bubbles/spinner`) for all async API calls.
- [ ] **Style Consolidation**: Move all hardcoded colors and paddings into `internal/ui/styles.go`.
- [ ] **Markdown Improvements**: Optimize `glamour` rendering to handle dark/light terminal themes dynamically.

## Phase 4: Testing & Quality Assurance ✅
- [ ] **Unit Testing**:
    - [ ] Achieve >80% coverage for `internal/api` (mock HTTP server).
    - [ ] Test TUI `Update` functions using `tea.Msg` simulation.
- [ ] **Integration Testing**: Implement golden file testing for TUI views to catch regression in layouts.
- [ ] **E2E Testing**:
    - [ ] Create a "Mock Ditto API" (or use a local containerized version).
    - [ ] Write CLI E2E tests using a tool like `expect` or a custom Go runner.
- [ ] **Linting**: Strict `golangci-lint` enforcement (currently failing due to version mismatch).

## Phase 5: CI/CD & DevEx 👷
- [ ] **GitHub Actions**:
    - [ ] Workflow for `lint` and `test` on every PR.
    - [ ] Automated release workflow (using `goreleaser`) for binaries.
- [ ] **Local Dev Environment**:
    - [ ] Add `docker-compose.yml` with a mock API and database for local testing.
    - [ ] Create a `scripts/seed.go` to populate the local API with test data.
- [ ] **Documentation**: Complete the `README.md` with a detailed troubleshooting guide and contribution instructions.

---
*Note: This list is dynamic and should be updated as the project evolves.*
