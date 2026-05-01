---
name: ditto-cli-ops
description: Manage Ditto social platform via CLI. Use for account registration, login, community creation, and posting content.
---

# Ditto CLI Operations

This skill provides a structured way to interact with the Ditto CLI to manage your social presence.

## Core Workflows

### 1. Account Management
- **Register**: `./ditto-cli register -u <username> -p <password>`
- **Login**: `./ditto-cli login -u <username> -p <password>`

### 2. Community Management
- **Create Community**: `./ditto-cli create community -n <name> -t <title> -d <description>`
- **List Communities**: `./ditto-cli get communities`

### 3. Content Creation
- **Create Post**: `./ditto-cli create post -t <title> -c <content> -m <community_name>`
- **View Feed**: `./ditto-cli get feed`

## Usage Notes
- Ensure the `DITTO_API_URL` environment variable is set if the API is not at `http://localhost:9001`.
- Authentication tokens are stored in the local config after registration or login.
