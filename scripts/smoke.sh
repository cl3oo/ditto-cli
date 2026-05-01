#!/usr/bin/env bash
set -euo pipefail

API_URL="${DITTO_API_URL:-http://localhost:9001/api/v1}"
BIN="${DITTO_SMOKE_BIN:-./dist/ditto-cli}"
STAMP="$(date +%s)"
RAND="${RANDOM}"
BASE="smoke-${STAMP}-${RAND}"
USER_ONE="${DITTO_SMOKE_USER_ONE:-${BASE}-owner}"
USER_TWO="${DITTO_SMOKE_USER_TWO:-${BASE}-member}"
PASS_ONE="${DITTO_SMOKE_PASS_ONE:-SmokePass123!}"
PASS_TWO="${DITTO_SMOKE_PASS_TWO:-SmokePass123!}"
EMAIL_ONE="${DITTO_SMOKE_EMAIL_ONE:-${USER_ONE}@example.com}"
EMAIL_TWO="${DITTO_SMOKE_EMAIL_TWO:-${USER_TWO}@example.com}"
COMMUNITY_NAME="${DITTO_SMOKE_COMMUNITY_NAME:-${BASE}-community}"
COMMUNITY_TITLE="${DITTO_SMOKE_COMMUNITY_TITLE:-${BASE} community}"
POST_TITLE="${DITTO_SMOKE_POST_TITLE:-${BASE}-post}"
POST_CONTENT="${DITTO_SMOKE_POST_CONTENT:-Product smoke test post from ditto-cli}"
TMP_HOME="$(mktemp -d)"

cleanup() {
  rm -rf "$TMP_HOME"
}
trap cleanup EXIT

if [[ ! -x "$BIN" ]]; then
  echo "Smoke binary not found or not executable: $BIN" >&2
  echo "Run 'make build' first or set DITTO_SMOKE_BIN." >&2
  exit 1
fi

run_cli() {
  HOME="$TMP_HOME" "$BIN" --url "$API_URL" "$@"
}

assert_contains() {
  local haystack="$1"
  local needle="$2"
  local label="$3"
  if ! grep -Fq "$needle" <<<"$haystack"; then
    echo "Expected $label to contain '$needle'" >&2
    echo "$haystack" >&2
    exit 1
  fi
}

echo "[1/9] Registering owner account"
owner_register_output="$(run_cli register -u "$USER_ONE" -e "$EMAIL_ONE" -p "$PASS_ONE")"
assert_contains "$owner_register_output" "Account created and logged in successfully" "register output"

echo "[2/9] Creating community $COMMUNITY_NAME"
owner_create_community_output="$(run_cli create community --name "$COMMUNITY_NAME" --title "$COMMUNITY_TITLE" --description "Smoke test community")"
assert_contains "$owner_create_community_output" "Community created successfully" "create community output"

echo "[3/9] Reading community list and resolving created community id"
communities_output="$(run_cli get communities)"
assert_contains "$communities_output" "$COMMUNITY_NAME" "communities output"
community_id="$(printf '%s\n' "$communities_output" | awk -v name="$COMMUNITY_NAME" '$2 == name { print $1; exit }')"
if [[ -z "$community_id" ]]; then
  echo "Could not resolve community id for $COMMUNITY_NAME" >&2
  echo "$communities_output" >&2
  exit 1
fi

echo "[4/9] Creating post $POST_TITLE"
owner_create_post_output="$(run_cli create post --title "$POST_TITLE" --content "$POST_CONTENT" --community "$COMMUNITY_NAME")"
assert_contains "$owner_create_post_output" "Post created successfully" "create post output"

echo "[5/9] Reading post listings and resolving created post id"
posts_output="$(run_cli get posts)"
assert_contains "$posts_output" "$POST_TITLE" "posts output"
post_id="$(printf '%s\n' "$posts_output" | awk -v title="$POST_TITLE" '$2 == title { print $1; exit }')"
if [[ -z "$post_id" ]]; then
  echo "Could not resolve post id for $POST_TITLE" >&2
  echo "$posts_output" >&2
  exit 1
fi

echo "[6/9] Reading personalized feed"
feed_output="$(run_cli get feed)"
if [[ -z "$feed_output" ]]; then
  echo "Feed output was empty" >&2
  exit 1
fi

echo "[7/9] Registering member account"
member_register_output="$(run_cli register -u "$USER_TWO" -e "$EMAIL_TWO" -p "$PASS_TWO")"
assert_contains "$member_register_output" "Account created and logged in successfully" "second register output"

echo "[8/9] Joining community as member"
member_join_output="$(run_cli join "$community_id")"
assert_contains "$member_join_output" "$community_id" "join output"

echo "[9/9] Voting on the created post as member"
member_vote_output="$(run_cli vote post "$post_id" --value 1)"
assert_contains "$member_vote_output" "$post_id" "vote output"

echo
echo "Smoke test passed"
echo "  community_id: $community_id"
echo "  post_id: $post_id"
echo "  owner: $USER_ONE"
echo "  member: $USER_TWO"
