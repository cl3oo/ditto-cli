# Implementation Plan & Technical Details: Ditto CLI

This document outlines the phased implementation of the remaining Ditto API features, including technical steps for the CLI application.

## Objective
To implement the missing API endpoints into the Ditto CLI, prioritizing core user operations, social features, moderation, and finally advanced ecosystem features.

---

## Phase 1: Core CRUD & Missing Endpoints

### Communities
- [x] Get all communities (`GET /communities`)
- [x] Get trending communities (`GET /communities/trending`)
- [x] Get community by ID (`GET /communities/{id}`)
- [x] Join/Leave community (`POST /communities/{id}/join`)
- [x] Create community (`POST /communities`)
- [x] **Update community** (`PUT /communities/{id}`)
  - **API**: Added `UpdateCommunity` in `internal/api/client.go`.
  - **TUI**: Added `:edit` command when viewing a community.
- [x] **Delete community** (`DELETE /communities/{id}`)
  - **API**: Added `DeleteCommunity` in `internal/api/client.go`.
  - **TUI**: Added `:delete` command when viewing a community.
- [x] **Get random communities** (`GET /communities/random`)
  - **API**: Added `GetRandomCommunities` in `internal/api/client.go`.
  - **TUI**: Added `:random` command in `StateCommunities`.

### Posts
- [x] Get all posts (filtered) (`GET /posts`)
- [x] Get trending posts (`GET /posts/trending`)
- [x] Get post by ID (`GET /posts/{id}`)
- [x] Create post (`POST /posts/{id}`)
- [x] **Update post** (`PUT /posts/{id}`)
  - **API**: Added `UpdatePost` in `internal/api/client.go`.
  - **TUI**: Added `:edit` command in `StatePostDetail`.
- [x] **Delete post** (`DELETE /posts/{id}`)
  - **API**: Added `DeletePost` in `internal/api/client.go`.
  - **TUI**: Added `:delete` command in `StatePostDetail`.
- [x] **Get random posts** (`GET /posts/random`)
  - **API**: Added `GetRandomPosts` in `internal/api/client.go`.
  - **TUI**: Added `:random` command in `StateFeed`.

### Comments
- [x] Get all comments (by target) (`GET /comments`)
- [x] Create comment (`POST /comments/{id}`)
- [x] **Delete comment** (`DELETE /comments/{id}`)
  - **API**: Added `DeleteComment` in `internal/api/client.go`.
  - **TUI**: Added `:delete-comment <id>` command.
- [x] **Get random comments** (`GET /comments/random`)
  - **API**: Added `GetRandomComments` in `internal/api/client.go`.
  - **CLI**: Added `ditto get comments --random 5`.

### Users
- [x] Get current user (`GET /users/me`)
- [x] Get user by ID (`GET /users/{id}`)
- [x] **Update user** (`PUT /users/{id}`)
  - **API**: Added `UpdateUser` in `internal/api/client.go`.
  - **TUI**: Added `StateProfileSettings` view via `:settings`.
- [x] **Delete user** (`DELETE /users/{id}`)
  - **API**: Added `DeleteUser` in `internal/api/client.go`.
  - **TUI/CLI**: Added `:delete-account` command.
- [x] **Get random users** (`GET /users/random`)
  - **API**: Added `GetRandomUsers` in `internal/api/client.go`.
  - **CLI**: Added `ditto get users --random 5`.

---

## Phase 2: Social & Subscription Features

### Connections
- [x] **Follow/Unfollow user** (`POST /users/{id}/follow`)
  - **API**: Added `ToggleFollow` in `internal/api/client.go`.
  - **TUI**: Added `:follow` command in `StatePostDetail`.
- [x] **Check following status** (`GET /users/{id}/check`)
  - **API**: Added `CheckFollowStatus` in `internal/api/client.go`.
  - **TUI**: Used by future profile views.
- [x] **Get followed users** (`GET /users/subed`)
  - **API**: Added `GetFollowedUsers` in `internal/api/client.go`.
  - **TUI**: Added `:following` command.

### Subscriptions
- [x] **Get joined communities** (`GET /communities/subed`)
  - **API**: Added `GetJoinedCommunities` in `internal/api/client.go`.
  - **TUI**: Added `:joined` command.
- [x] **Check if joined community** (`GET /communities/{id}/check`)
  - **API**: Added `CheckCommunityStatus` in `internal/api/client.go`.

---

## Phase 3: Moderation Tools

### Community Moderation
- [x] **Ban a user** (`POST /communities/{id}/ban`)
  - **API**: Added `BanUser` in `internal/api/client.go`.
  - **TUI**: Added `:ban <userID>` command.
- [x] **Add a moderator** (`POST /communities/{id}/mods`)
  - **API**: Added `AddModerator` in `internal/api/client.go`.
  - **TUI**: Added `:mod add <userID>` command.
- [x] **Update moderator permissions** (`PUT /communities/{id}/mods/{userId}`)
  - **API**: Added `UpdateModerator` in `internal/api/client.go`.
- [x] **Remove a moderator** (`DELETE /communities/{id}/mods/{userId}`)
  - **API**: Added `RemoveModerator` in `internal/api/client.go`.

### Post Moderation
- [x] **Lock a post** (`POST /posts/{id}/lock`)
  - **API**: Added `LockPost` in `internal/api/client.go`.
  - **TUI**: Added `:lock` and `:unlock` commands. Comment input is hidden if locked.
- [x] **Moderator delete post** (`DELETE /posts/{id}/mod`)
  - **API**: Added `ModDeletePost` in `internal/api/client.go`.
  - **TUI**: Added `:mod-delete <reason>` command.

---

## Phase 4: Advanced Features & Ecosystem

### Media
- [x] **Upload media** (`POST /media/{id}`)
  - **API**: Added `UploadMedia` in `internal/api/client.go`.
  - **CLI**: Added `ditto upload <file> --target <id>`.
- [x] **Serve media** (`GET /media/{id}`)
  - **API**: Added `DownloadMedia` in `internal/api/client.go`.
  - **CLI**: Added `ditto download <id> --output <path>`.

### Wallet & Economy
- [x] **Get current user's wallet** (`GET /users/wallet`)
  - **API**: Added `GetWallet` in `internal/api/client.go`.
  - **TUI**: Displaying coins and tokens in the header.
- [x] **Give an award** (`POST /awards/{id}`)
  - **API**: Added `GiveAward` in `internal/api/client.go`.
  - **TUI**: Added `:award <award_id>` command.
- [x] **Remove an award** (`DELETE /awards/{id}`)
  - **API**: Added `RemoveAward` in `internal/api/client.go`.

### Platform Operations
- [x] **Submit a report** (`POST /reports/{id}`)
  - **API**: Added `ReportResource` in `internal/api/client.go`.
  - **TUI**: Added `:report <reason>` command.
- [x] **Get resource metadata** (`GET /meta/{id}`)
  - **API**: Added `GetMetaData` in `internal/api/client.go`.
  - **TUI**: Used for background updates and refresh.

---

## Completed Base Functionality
- [x] Authorize a user (`POST /auth/authorize`)
- [x] Register a user (`POST /auth/register`)
- [x] Get user's home feed (`GET /feed`)
- [x] Submit a vote (Post/Comment) (`POST /votes/{id}`)
- [x] Search communities (`GET /search`)
