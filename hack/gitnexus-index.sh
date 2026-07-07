#!/usr/bin/env bash

#   /    ctx:                         https://ctx.ist
# ,'`./    do you remember?
# `.,'\
#   \    Copyright 2026-present Context contributors.
#                 SPDX-License-Identifier: Apache-2.0
#
# gitnexus-index.sh — index this repo into GitNexus via gitnexus-docker.sh,
# and (optionally) register the MCP server with Claude Code.
#
# Vendored from the sibling os repo's scripts/ directory.
#
#   ./hack/gitnexus-index.sh                 # index this repo
#   REGISTER_MCP=1 ./hack/gitnexus-index.sh  # also `claude mcp add` the server
#   ./hack/gitnexus-index.sh <repo> ...      # index specific repos instead
#
# The index lands in <repo>/.gitnexus/ (gitignored). GitNexus's analyze needs
# network for the one-time embedding-model fetch; everything else is offline.
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
WRAPPER="$SCRIPT_DIR/gitnexus-docker.sh"
REGISTER_MCP="${REGISTER_MCP:-0}"

# Default target: the repo this script is vendored into. Override by passing
# one or more repo paths as arguments.
if [ "$#" -gt 0 ]; then
  REPOS=("$@")
else
  REPOS=("$(git -C "$SCRIPT_DIR" rev-parse --show-toplevel)")
fi

# ---- pre-flight -----------------------------------------------------------
command -v docker >/dev/null 2>&1 || { echo "docker not found on PATH" >&2; exit 1; }
[ -x "$WRAPPER" ] || { echo "missing or non-executable: $WRAPPER" >&2; exit 1; }

# ---- index each repo ------------------------------------------------------
indexed=0
for repo in "${REPOS[@]}"; do
  echo
  echo "==================================================================="
  echo "  $repo"
  echo "==================================================================="
  if [ ! -d "$repo" ]; then echo "  SKIP — not a directory"; continue; fi
  if ! git -C "$repo" rev-parse --git-dir >/dev/null 2>&1; then
    echo "  SKIP — not a git repo (GitNexus needs a git root)"; continue
  fi

  # keep GitNexus's index out of the repo's git history
  gi="$repo/.gitignore"
  if [ -f "$gi" ] && grep -qxF '.gitnexus/' "$gi"; then
    echo "  '.gitnexus/' already gitignored"
  else
    printf '\n# GitNexus code-graph index (local, do not commit)\n.gitnexus/\n' >> "$gi"
    echo "  added '.gitnexus/' to $gi"
  fi

  if "$WRAPPER" index "$repo"; then
    indexed=$((indexed + 1))
  else
    echo "  !! indexing FAILED for $repo (continuing)"
  fi
done

echo
echo "=== registered repositories ==="
"$WRAPPER" list || true
echo
echo ">> indexed $indexed repo(s)."

# ---- optional: register the MCP server with Claude Code -------------------
MCP_CMD=(claude mcp add --scope user gitnexus -- "$WRAPPER" mcp)
echo
if [ "$REGISTER_MCP" = "1" ]; then
  echo ">> registering the GitNexus MCP server with Claude Code"
  "${MCP_CMD[@]}"
  echo ">> done — RESTART Claude Code for it to connect."
else
  echo "To wire the MCP server into Claude Code, re-run with REGISTER_MCP=1,"
  echo "or run this yourself (then restart Claude Code):"
  printf '   '; printf '%q ' "${MCP_CMD[@]}"; echo
fi
