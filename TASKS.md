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
- [ ] **Update community** (`PUT /communities/{id}`)
  - **API**: Add `UpdateCommunity(id string, data map[string]interface{}) error` in `internal/api/client.go`.
  - **TUI**: Add `:edit` command when viewing a community. Create an `EditCommunityModel` similar to `CreatePostModel` to handle input.
- [ ] **Delete community** (`DELETE /communities/{id}`)
  - **API**: Add `DeleteCommunity(id string) error` in `internal/api/client.go`.
  - **TUI**: Add `:delete` command when viewing a community. Prompt for confirmation (`y/n`).
- [ ] **Get random communities** (`GET /communities/random`)
  - **API**: Add `GetRandomCommunities(num int) ([]types.Community, error)` in `internal/api/client.go`.
  - **TUI**: Add `:random` command in `StateCommunities` to trigger this fetch and update the list.

### Posts
- [x] Get all posts (filtered) (`GET /posts`)
- [x] Get trending posts (`GET /posts/trending`)
- [x] Get post by ID (`GET /posts/{id}`)
- [x] Create post (`POST /posts/{id}`)
- [ ] **Update post** (`PUT /posts/{id}`)
  - **API**: Add `UpdatePost(id, title, content string) error` in `internal/api/client.go`.
  - **TUI**: Add `:edit` command in `StatePostDetail`. Create an `EditPostModel` allowing the user to modify the title/content.
- [ ] **Delete post** (`DELETE /posts/{id}`)
  - **API**: Add `DeletePost(id string) error` in `internal/api/client.go`.
  - **TUI**: Add `:delete` command in `StatePostDetail`. Prompt for confirmation and return to `StateFeed` on success.
- [ ] **Get random posts** (`GET /posts/random`)
  - **API**: Add `GetRandomPosts(num int) ([]types.Post, error)` in `internal/api/client.go`.
  - **TUI**: Add `:random` command in `StateFeed` to populate the feed with random posts.

### Comments
- [x] Get all comments (by target) (`GET /comments`)
- [x] Create comment (`POST /comments/{id}`)
- [ ] **Delete comment** (`DELETE /comments/{id}`)
  - **API**: Add `DeleteComment(id string) error` in `internal/api/client.go`.
  - **TUI**: Add keybinding (e.g., `D` or `:delete-comment <id>`) in `StatePostDetail` to remove a specific comment.
- [ ] **Get random comments** (`GET /comments/random`)
  - **API**: Add `GetRandomComments(num int) ([]types.Comment, error)` in `internal/api/client.go`.
  - **CLI**: Add a specific CLI command `ditto get comments --random 5`.

### Users
- [x] Get current user (`GET /users/me`)
- [x] Get user by ID (`GET /users/{id}`)
- [ ] **Update user** (`PUT /users/{id}`)
  - **API**: Add `UpdateUser(id string, data map[string]interface{}) error` in `internal/api/client.go`.
  - **TUI**: Add a `StateProfileSettings` view to edit avatar, bio, or other editable user fields.
- [ ] **Delete user** (`DELETE /users/{id}`)
  - **API**: Add `DeleteUser(id string) error` in `internal/api/client.go`.
  - **TUI/CLI**: Add extreme confirmation prompt before calling this endpoint (account deletion).
- [ ] **Get random users** (`GET /users/random`)
  - **API**: Add `GetRandomUsers(num int) ([]types.User, error)` in `internal/api/client.go`.

---

## Phase 2: Social & Subscription Features

### Connections
- [ ] **Follow/Unfollow user** (`POST /users/{id}/follow`)
  - **API**: Add `ToggleFollow(userID string) error` in `internal/api/client.go`.
  - **TUI**: Add `:follow` command when viewing a user profile or post author.
- [ ] **Check following status** (`GET /users/{id}/check`)
  - **API**: Add `CheckFollowStatus(userID string) (bool, error)` in `internal/api/client.go`.
  - **TUI**: Call this when loading a user profile to show `[Following]` or `[Follow]` button.
- [ ] **Get followed users** (`GET /users/subed`)
  - **API**: Add `GetFollowedUsers() ([]types.User, error)` in `internal/api/client.go`.
  - **TUI**: Add a specific view or a `:following` command to list these users.

### Subscriptions
- [ ] **Get joined communities** (`GET /communities/subed`)
  - **API**: Add `GetJoinedCommunities() ([]types.Community, error)` in `internal/api/client.go`.
  - **TUI**: Add a `:joined` command in `StateCommunities` to filter the list to only subscribed ones.
- [ ] **Check if joined community** (`GET /communities/{id}/check`)
  - **API**: Add `CheckCommunityStatus(id string) (bool, error)` in `internal/api/client.go`.
  - **TUI**: Ensure the community detail view reflects the correct `Join/Leave` state on load.

---

## Phase 3: Moderation Tools

### Community Moderation
- [ ] **Ban a user** (`POST /communities/{id}/ban`)
  - **API**: Add `BanUser(communityID, userID string) error` in `internal/api/client.go`.
  - **TUI**: Add `:ban <username>` command in the community view (requires resolving username to ID or adding search).
- [ ] **Add a moderator** (`POST /communities/{id}/mods`)
  - **API**: Add `AddModerator(communityID, userID string) error` in `internal/api/client.go`.
  - **TUI**: Add `:mod add <username>` command.
- [ ] **Update moderator permissions** (`PUT /communities/{id}/mods/{userId}`)
  - **API**: Add `UpdateModerator(communityID, userID string, permissions map[string]interface{}) error` in `internal/api/client.go`.
- [ ] **Remove a moderator** (`DELETE /communities/{id}/mods/{userId}`)
  - **API**: Add `RemoveModerator(communityID, userID string) error` in `internal/api/client.go`.

### Post Moderation
- [ ] **Lock a post** (`POST /posts/{id}/lock`)
  - **API**: Add `LockPost(postID string, lock bool) error` in `internal/api/client.go`.
  - **TUI**: Add `:lock` command in `StatePostDetail`. Hide comment input if post is locked.
- [ ] **Moderator delete post** (`DELETE /posts/{id}/mod`)
  - **API**: Add `ModDeletePost(postID, reason string) error` in `internal/api/client.go`.
  - **TUI**: Add `:mod-delete <reason>` command in `StatePostDetail`.

---

## Phase 4: Advanced Features & Ecosystem

### Media
- [ ] **Upload media** (`POST /media/{id}`)
  - **API**: Add `UploadMedia(targetID string, targetType int, filePath string) error` using `multipart/form-data` in `internal/api/client.go`.
  - **CLI**: Implement a specific CLI command `ditto upload <file> --target <id>`.
- [ ] **Serve media** (`GET /media/{id}`)
  - **API**: Add `DownloadMedia(mediaID, outputPath string) error` in `internal/api/client.go`.

### Wallet & Economy
- [ ] **Get current user's wallet** (`GET /users/wallet`)
  - **API**: Add `GetWallet() (*types.Wallet, error)` in `internal/api/client.go`.
  - **TUI**: Display token/coin balance in the main header or profile view.
- [ ] **Give an award** (`POST /awards/{id}`)
  - **API**: Add `GiveAward(targetID string, targetType int, awardID string) error` in `internal/api/client.go`.
  - **TUI**: Add `:award <type>` command in `StatePostDetail`.
- [ ] **Remove an award** (`DELETE /awards/{id}`)
  - **API**: Add `RemoveAward(awardID string) error` in `internal/api/client.go`.

### Platform Operations
- [ ] **Submit a report** (`POST /reports/{id}`)
  - **API**: Add `ReportResource(targetID string, targetType int, reason string) error` in `internal/api/client.go`.
  - **TUI**: Add `:report <reason>` command when viewing posts or comments.
- [ ] **Get resource metadata** (`GET /meta/{id}`)
  - **API**: Add `GetMetaData(targetID string, targetType int) (*types.MetaDataResult, error)` in `internal/api/client.go`.
  - **TUI**: Call this periodically or on-load to update vote scores and award counts.

---

## Completed Base Functionality
- [x] Authorize a user (`POST /auth/authorize`)
- [x] Register a user (`POST /auth/register`)
- [x] Get user's home feed (`GET /feed`)
- [x] Submit a vote (Post/Comment) (`POST /votes/{id}`)
- [x] Search communities (`GET /search`)
